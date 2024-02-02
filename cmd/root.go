/*
Copyright © 2023 Stephan Bechter <s.bechter@netconomy.net>
*/
package cmd

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/grid-x/modbus"
	"github.com/spf13/cobra"
)

const connectTimeout = 2 * time.Second

var rootCmdArgs struct {
	ip     string
	port   int
	unitID int
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "huawei-check",
	Short: "Check the ModbusTCP communication to huawei inverters",
	Run: func(cmd *cobra.Command, args []string) {
		handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", rootCmdArgs.ip, rootCmdArgs.port))
		handler.Timeout = 10 * time.Second
		handler.SlaveID = byte(rootCmdArgs.unitID)
		handler.ConnectDelay = connectTimeout

		fmt.Printf("Connecting to inverter... (Waiting %s after connected)\n", connectTimeout)
		err := handler.Connect()
		cobra.CheckErr(err)
		defer handler.Close()

		client := modbus.NewClient(handler)

		results, err := client.ReadHoldingRegisters(30000, 15)
		cobra.CheckErr(err)
		fmt.Printf("Model        : [%s]\n", string(results))

		results, err = client.ReadHoldingRegisters(30070, 1)
		cobra.CheckErr(err)
		fmt.Printf("Model ID     : [%d]\n", int16(binary.BigEndian.Uint16(results)))

		results, err = client.ReadHoldingRegisters(30015, 10)
		cobra.CheckErr(err)
		fmt.Printf("Serial number: [%s]\n", string(results))

		results, err = client.ReadHoldingRegisters(30073, 2)
		cobra.CheckErr(err)
		fmt.Printf("Rated power  : [%d W]\n", int32(binary.BigEndian.Uint32(results)))

		results, err = client.ReadHoldingRegisters(32064, 2)
		cobra.CheckErr(err)
		fmt.Printf("Input power  : [%d W]\n", int32(binary.BigEndian.Uint32(results)))

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootCmdArgs.ip, "ip", "", "Inverter IP address")
	rootCmd.PersistentFlags().IntVar(&rootCmdArgs.port, "port", 502, "Inverter ModbusTCP port")
	rootCmd.PersistentFlags().IntVar(&rootCmdArgs.unitID, "unitID", 1, "Modbus UnitID")

	_ = rootCmd.MarkPersistentFlagRequired("ip")
}
