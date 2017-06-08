package utils

// global
var Version = byte('\x05')
var ReplySuccess = byte('\x00')
var Rsv = byte('\x00')

// auth
var AuthNo = byte('\x00')
var AuthUsernamePassword = byte('\x02')
var NoAcceptableMethods = byte('\xFF')

// command
var CmdConnect = byte('\x01')

// ATYP
var AptyIPV4 = byte('\x01')
var AptyDomainName = byte('\x03')
var AptyIPV6 = byte('\x04')
