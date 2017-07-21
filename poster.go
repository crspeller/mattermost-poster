package main

import (
	"bytes"
	"io"
	"os"

	"fmt"

	"github.com/howeyc/gopass"
	"github.com/mattermost/platform/model"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mattermost-poster [server]",
	Short: "Post messages with attachments to Mattermost",
	RunE:  doPostCmdF,
}

func main() {
	rootCmd.Flags().StringP("username", "u", "", "Username to login with")
	rootCmd.Flags().StringP("password", "p", "", "Password to login with")
	rootCmd.Flags().StringP("message", "m", "", "Text to send")
	rootCmd.Flags().StringP("channel", "c", "", "The channel ID to send message to")
	//rootCmd.Flags().StringP("fmessage", "f", "", "File to send as a message")
	rootCmd.Flags().StringArrayP("attachment", "a", []string{}, "File to attach")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func doPostCmdF(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Need a server URL")
	}

	if len(args) > 1 {
		return fmt.Errorf("Extra args")
	}

	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	message, _ := cmd.Flags().GetString("message")
	channelId, _ := cmd.Flags().GetString("channel")
	attachments, _ := cmd.Flags().GetStringArray("attachment")

	if password == "" {
		fmt.Print("Password: ")
		getpass, err := gopass.GetPasswd()
		if err != nil {
			return fmt.Errorf("Need a password")
		}
		password = string(getpass)
	}

	client := model.NewAPIv4Client(args[0])

	user, resp := client.Login(username, password)
	if resp.Error != nil {
		return resp.Error
	}

	var fileIds []string
	if len(attachments) != 0 {
		for _, filename := range attachments {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Print("Unable to find: " + filename)
				fmt.Println(" Error: " + err.Error())
				continue
			}
			data := &bytes.Buffer{}
			if _, err := io.Copy(data, file); err != nil {
				fmt.Print("Unable to copy file: " + filename)
				fmt.Println(" Error: " + err.Error())
				continue
			}
			file.Close()

			fileUploadResp, resp := client.UploadFile(data.Bytes(), channelId, filename)
			if resp.Error != nil || fileUploadResp == nil || len(fileUploadResp.FileInfos) != 1 {
				fmt.Print("Unable to upload file: " + filename)
				fmt.Println(" Error: " + resp.Error.Error())
				continue
			}

			fileIds = append(fileIds, fileUploadResp.FileInfos[0].Id)
		}
	}

	client.CreatePost(&model.Post{
		UserId:    user.Id,
		ChannelId: channelId,
		Message:   message,
		Type:      model.POST_DEFAULT,
		FileIds:   fileIds,
	})

	return nil
}
