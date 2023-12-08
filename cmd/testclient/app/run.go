package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Joffref/slotframe-manager/internal/api"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/udp"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var (
	slotframeVersion uint16
	slots            api.Slots
	runCmd           = &cobra.Command{
		Use:   "run",
		Short: "Run the slotframe-manager testclient",
		Long:  "Run the slotframe-manager testclient",
		RunE: func(cmd *cobra.Command, args []string) error {
			co, err := udp.Dial("localhost:5688")
			if err != nil {
				log.Fatalf("Error dialing: %v", err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Get version of slotframe
			url := "/version"
			resp, err := co.Get(ctx, url)
			if err != nil {
				log.Fatalf("Error sending request: %v", err)
			}
			fmt.Println(resp.String())
			var v *api.FrameVersion
			err = json.NewDecoder(resp.Body()).Decode(&v)
			if err != nil {
				return err
			}
			log.Printf("Response: %v", v)

			fmt.Println(parentId, id, etx)

			url = fmt.Sprintf("/register/%v/%v/%v", parentId, id, etx)

			body, err := json.Marshal(api.Slots{})
			if err != nil {
				log.Fatalf("Error marshalling body: %v", err)
			}
			fmt.Println(string(url))
			resp, err = co.Post(ctx, url, message.AppJSON, bytes.NewReader(body))
			if err != nil {
				log.Fatalf("Error sending request: %v", err)
			}
			fmt.Println(resp.String())
			slots = api.Slots{}
			err = json.NewDecoder(resp.Body()).Decode(&slots)
			if err != nil {
				log.Fatalf("Error decoding response: %v", err)
			}
			var slots api.Slots
			var version api.FrameVersion
			for {
				time.Sleep(1 * time.Second)
				resp, err = co.Get(ctx, "/version")
				if err != nil {
					log.Fatalf("Error sending request: %v", err)
				}
				err = json.NewDecoder(resp.Body()).Decode(&version)
				if err != nil {
					return err
				}
				if uint16(version.Version) != slotframeVersion {
					log.Printf("Version mismatch: %v != %v", v, slotframeVersion)
					// retrieve slots
					url = fmt.Sprintf("/register/%v/%v/%v", parentId, id, etx)
					body, err := json.Marshal(slots)
					if err != nil {
						log.Fatalf("Error marshalling body: %v", err)
					}
					resp, err = co.Post(ctx, url, message.AppJSON, bytes.NewReader(body))
					if err != nil {
						log.Fatalf("Error sending request: %v", err)
					}
					err = json.NewDecoder(resp.Body()).Decode(&slots)
					if err != nil {
						return err
					}
					log.Printf("Response: %v", slots)
					slotframeVersion = uint16(version.Version)
				}
			}
		},
	}
)
