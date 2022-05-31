package ssh

import "golang.org/x/crypto/ssh"

//指定SSH-2允许使用的加密算法。多个算法之间使用逗号分隔。可以使用的算法如下：
//"aes128-cbc", "aes192-cbc", "aes256-cbc", "aes128-ctr", "aes192-ctr", "aes256-ctr",
//"3des-cbc", "arcfour128", "arcfour256", "arcfour", "blowfish-cbc", "cast128-cbc"
//默认值是可以使用上述所有算法。
func GetSshSupportedCiphers() []string {
	config := &ssh.ClientConfig{}
	config.SetDefaults()
	for _, cipher := range []string{"aes128-cbc"} {
		found := false
		for _, defaultCipher := range config.Ciphers {
			if cipher == defaultCipher {
				found = true
				break
			}
		}

		if !found {
			config.Ciphers = append(config.Ciphers, cipher)
		}
	}
	return config.Ciphers
}
