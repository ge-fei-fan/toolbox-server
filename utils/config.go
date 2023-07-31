package utils

//
//import (
//	"bytes"
//	"encoding/json"
//	"toolbox-server/global"
//)
//
//func WriteConfig() error {
//	config, err := json.Marshal(global.TOOL_CONFIG)
//	if err != nil {
//		return err
//	}
//	err = global.TOOL_VP.ReadConfig(bytes.NewBuffer(config))
//	if err != nil {
//		return err
//	}
//	err = global.TOOL_VP.WriteConfig()
//	if err != nil {
//		return err
//	}
//	return nil
//}
