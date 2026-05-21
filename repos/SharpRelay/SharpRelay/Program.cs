/*
* Copyright © 2018-Present Jesse Nicholson
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Runtime.InteropServices;
using System.Threading;
using WinDivertSharp;
using WinDivertSharp.Extensions;
using WinDivertSharp.WinAPI;

namespace SharpRelay
{
    public class Program
    {

        // https://gist.githubusercontent.com/FusRoDah061/d04dc0bbed890ba0e93166da2b62451e/raw/47a0b27c86f25b9c3291197dd2f6dd620450d61e/ServiceInstaller.cs
        
        public static class ServiceInstaller
        {
            private const int STANDARD_RIGHTS_REQUIRED = 0xF0000;
            private const int SERVICE_WIN32_OWN_PROCESS = 0x00000010;

            [StructLayout(LayoutKind.Sequential)]
            private class SERVICE_STATUS
            {
                public int dwServiceType = 0;
                public ServiceState dwCurrentState = 0;
                public int dwControlsAccepted = 0;
                public int dwWin32ExitCode = 0;
                public int dwServiceSpecificExitCode = 0;
                public int dwCheckPoint = 0;
                public int dwWaitHint = 0;
            }

            #region OpenSCManager
            [DllImport("advapi32.dll", EntryPoint = "OpenSCManagerW", ExactSpelling = true, CharSet = CharSet.Unicode, SetLastError = true)]
            static extern IntPtr OpenSCManager(string machineName, string databaseName, ScmAccessRights dwDesiredAccess);
            #endregion

            #region OpenService
            [DllImport("advapi32.dll", SetLastError = true, CharSet = CharSet.Auto)]
            static extern IntPtr OpenService(IntPtr hSCManager, string lpServiceName, ServiceAccessRights dwDesiredAccess);
            #endregion

            #region CreateService
            [DllImport("advapi32.dll", SetLastError = true, CharSet = CharSet.Auto)]
            private static extern IntPtr CreateService(IntPtr hSCManager, string lpServiceName, string lpDisplayName, ServiceAccessRights dwDesiredAccess, int dwServiceType, ServiceBootFlag dwStartType, ServiceError dwErrorControl, string lpBinaryPathName, string lpLoadOrderGroup, IntPtr lpdwTagId, string lpDependencies, string lp, string lpPassword);
            #endregion

            #region CloseServiceHandle
            [DllImport("advapi32.dll", SetLastError = true)]
            [return: MarshalAs(UnmanagedType.Bool)]
            static extern bool CloseServiceHandle(IntPtr hSCObject);
            #endregion

            #region QueryServiceStatus
            [DllImport("advapi32.dll")]
            private static extern int QueryServiceStatus(IntPtr hService, SERVICE_STATUS lpServiceStatus);
            #endregion

            #region DeleteService
            [DllImport("advapi32.dll", SetLastError = true)]
            [return: MarshalAs(UnmanagedType.Bool)]
            private static extern bool DeleteService(IntPtr hService);
            #endregion

            #region ControlService
            [DllImport("advapi32.dll")]
            private static extern int ControlService(IntPtr hService, ServiceControl dwControl, SERVICE_STATUS lpServiceStatus);
            #endregion

            #region StartService
            [DllImport("advapi32.dll", SetLastError = true)]
            private static extern int StartService(IntPtr hService, int dwNumServiceArgs, int lpServiceArgVectors);
            #endregion

            public static void Uninstall(string serviceName)
            {
                IntPtr scm = OpenSCManager(ScmAccessRights.AllAccess);

                try
                {
                    IntPtr service = OpenService(scm, serviceName, ServiceAccessRights.AllAccess);
                    if (service == IntPtr.Zero)
                        throw new ApplicationException("Service not installed.");

                    try
                    {
                        StopService(service);
                        if (!DeleteService(service))
                            throw new ApplicationException("Could not delete service " + Marshal.GetLastWin32Error());
                    }
                    finally
                    {
                        CloseServiceHandle(service);
                    }
                }
                finally
                {
                    CloseServiceHandle(scm);
                }
            }

            public static bool ServiceIsInstalled(string serviceName)
            {
                IntPtr scm = OpenSCManager(ScmAccessRights.Connect);

                try
                {
                    IntPtr service = OpenService(scm, serviceName, ServiceAccessRights.QueryStatus);

                    if (service == IntPtr.Zero)
                        return false;

                    CloseServiceHandle(service);
                    return true;
                }
                finally
                {
                    CloseServiceHandle(scm);
                }
            }

            public static void InstallAndStart(string serviceName, string displayName, string fileName)
            {
                IntPtr scm = OpenSCManager(ScmAccessRights.AllAccess);

                try
                {
                    IntPtr service = OpenService(scm, serviceName, ServiceAccessRights.AllAccess);

                    if (service == IntPtr.Zero)
                        service = CreateService(scm, serviceName, displayName, ServiceAccessRights.AllAccess, 0x00000001, ServiceBootFlag.DemandStart, ServiceError.Normal, fileName, null, IntPtr.Zero, null, null, null);

                    if (service == IntPtr.Zero)
                        throw new ApplicationException("Failed to install service.");

                    try
                    {
                        StartService(service);
                    }
                    finally
                    {
                        CloseServiceHandle(service);
                    }
                }
                finally
                {
                    CloseServiceHandle(scm);
                }
            }

            public static void Install(string serviceName, string displayName, string fileName)
            {
                IntPtr scm = OpenSCManager(ScmAccessRights.AllAccess);

                try
                {
                    IntPtr service = OpenService(scm, serviceName, ServiceAccessRights.AllAccess);

                    if (service == IntPtr.Zero)
                        service = CreateService(scm, serviceName, displayName, ServiceAccessRights.AllAccess, SERVICE_WIN32_OWN_PROCESS, ServiceBootFlag.AutoStart, ServiceError.Normal, fileName, null, IntPtr.Zero, null, null, null);

                    if (service == IntPtr.Zero)
                        throw new ApplicationException("Failed to install service.");
                }
                finally
                {
                    CloseServiceHandle(scm);
                }
            }

            public static void StartService(string serviceName)
            {
                IntPtr scm = OpenSCManager(ScmAccessRights.Connect);

                try
                {
                    IntPtr service = OpenService(scm, serviceName, ServiceAccessRights.QueryStatus | ServiceAccessRights.Start);
                    if (service == IntPtr.Zero)
                        throw new ApplicationException("Could not open service.");

                    try
                    {
                        StartService(service);
                    }
                    finally
                    {
                        CloseServiceHandle(service);
                    }
                }
                finally
                {
                    CloseServiceHandle(scm);
                }
            }

            public static void StopService(string serviceName)
            {
                IntPtr scm = OpenSCManager(ScmAccessRights.Connect);

                try
                {
                    IntPtr service = OpenService(scm, serviceName, ServiceAccessRights.QueryStatus | ServiceAccessRights.Stop);
                    if (service == IntPtr.Zero)
                        throw new ApplicationException("Could not open service.");

                    try
                    {
                        StopService(service);
                    }
                    finally
                    {
                        CloseServiceHandle(service);
                    }
                }
                finally
                {
                    CloseServiceHandle(scm);
                }
            }

            private static void StartService(IntPtr service)
            {
                SERVICE_STATUS status = new SERVICE_STATUS();
                StartService(service, 0, 0);
                var changedStatus = WaitForServiceStatus(service, ServiceState.StartPending, ServiceState.Running);
                if (!changedStatus)
                    throw new ApplicationException("Unable to start service");
            }

            private static void StopService(IntPtr service)
            {
                SERVICE_STATUS status = new SERVICE_STATUS();
                ControlService(service, ServiceControl.Stop, status);
                var changedStatus = WaitForServiceStatus(service, ServiceState.StopPending, ServiceState.Stopped);
                if (!changedStatus)
                    throw new ApplicationException("Unable to stop service");
            }

            public static ServiceState GetServiceStatus(string serviceName)
            {
                IntPtr scm = OpenSCManager(ScmAccessRights.Connect);

                try
                {
                    IntPtr service = OpenService(scm, serviceName, ServiceAccessRights.QueryStatus);
                    if (service == IntPtr.Zero)
                        return ServiceState.NotFound;

                    try
                    {
                        return GetServiceStatus(service);
                    }
                    finally
                    {
                        CloseServiceHandle(service);
                    }
                }
                finally
                {
                    CloseServiceHandle(scm);
                }
            }

            private static ServiceState GetServiceStatus(IntPtr service)
            {
                SERVICE_STATUS status = new SERVICE_STATUS();

                if (QueryServiceStatus(service, status) == 0)
                    throw new ApplicationException("Failed to query service status.");

                return status.dwCurrentState;
            }

            private static bool WaitForServiceStatus(IntPtr service, ServiceState waitStatus, ServiceState desiredStatus)
            {
                SERVICE_STATUS status = new SERVICE_STATUS();

                QueryServiceStatus(service, status);
                if (status.dwCurrentState == desiredStatus) return true;

                int dwStartTickCount = Environment.TickCount;
                int dwOldCheckPoint = status.dwCheckPoint;

                while (status.dwCurrentState == waitStatus)
                {
                    // Do not wait longer than the wait hint. A good interval is
                    // one tenth the wait hint, but no less than 1 second and no
                    // more than 10 seconds.

                    int dwWaitTime = status.dwWaitHint / 10;

                    if (dwWaitTime < 1000) dwWaitTime = 1000;
                    else if (dwWaitTime > 10000) dwWaitTime = 10000;

                    Thread.Sleep(dwWaitTime);

                    // Check the status again.

                    if (QueryServiceStatus(service, status) == 0) break;

                    if (status.dwCheckPoint > dwOldCheckPoint)
                    {
                        // The service is making progress.
                        dwStartTickCount = Environment.TickCount;
                        dwOldCheckPoint = status.dwCheckPoint;
                    }
                    else
                    {
                        if (Environment.TickCount - dwStartTickCount > status.dwWaitHint)
                        {
                            // No progress made within the wait hint
                            break;
                        }
                    }
                }
                return (status.dwCurrentState == desiredStatus);
            }

            private static IntPtr OpenSCManager(ScmAccessRights rights)
            {
                IntPtr scm = OpenSCManager(null, null, rights);
                if (scm == IntPtr.Zero)
                    throw new ApplicationException("Could not connect to service control manager.");

                return scm;
            }
        }


        public enum ServiceState
        {
            Unknown = -1, // The state cannot be (has not been) retrieved.
            NotFound = 0, // The service is not known on the host server.
            Stopped = 1,
            StartPending = 2,
            StopPending = 3,
            Running = 4,
            ContinuePending = 5,
            PausePending = 6,
            Paused = 7
        }

        [Flags]
        public enum ScmAccessRights
        {
            Connect = 0x0001,
            CreateService = 0x0002,
            EnumerateService = 0x0004,
            Lock = 0x0008,
            QueryLockStatus = 0x0010,
            ModifyBootConfig = 0x0020,
            StandardRightsRequired = 0xF0000,
            AllAccess = (StandardRightsRequired | Connect | CreateService |
                         EnumerateService | Lock | QueryLockStatus | ModifyBootConfig)
        }

        [Flags]
        public enum ServiceAccessRights
        {
            QueryConfig = 0x1,
            ChangeConfig = 0x2,
            QueryStatus = 0x4,
            EnumerateDependants = 0x8,
            Start = 0x10,
            Stop = 0x20,
            PauseContinue = 0x40,
            Interrogate = 0x80,
            UserDefinedControl = 0x100,
            Delete = 0x00010000,
            StandardRightsRequired = 0xF0000,
            AllAccess = (StandardRightsRequired | QueryConfig | ChangeConfig |
                         QueryStatus | EnumerateDependants | Start | Stop | PauseContinue |
                         Interrogate | UserDefinedControl)
        }

        public enum ServiceBootFlag
        {
            Start = 0x00000000,
            SystemStart = 0x00000001,
            AutoStart = 0x00000002,
            DemandStart = 0x00000003,
            Disabled = 0x00000004
        }

        public enum ServiceControl
        {
            Stop = 0x00000001,
            Pause = 0x00000002,
            Continue = 0x00000003,
            Interrogate = 0x00000004,
            Shutdown = 0x00000005,
            ParamChange = 0x00000006,
            NetBindAdd = 0x00000007,
            NetBindRemove = 0x00000008,
            NetBindEnable = 0x00000009,
            NetBindDisable = 0x0000000A
        }

        public enum ServiceError
        {
            Ignore = 0x00000000,
            Normal = 0x00000001,
            Severe = 0x00000002,
            Critical = 0x00000003
        }


        public static ushort _originalPort;
        public static ushort _modifiedPort;

        public static volatile bool s_running = true;
        public static void Main(string[] args)
        {
            if (args.Length == 2)
            {
                if (args[0] == "uninstall")
                {
                    if (ServiceInstaller.ServiceIsInstalled(args[1]))
                    {
                        Console.WriteLine("[-] Removing service for loading driver...");
                        ServiceInstaller.Uninstall(args[1]);
                        Console.WriteLine("[-] Service " + args[1] + " was removed successfully!");
                        return;
                    }
                    else
                    {
                        Console.WriteLine("[!] Service " + args[1] + " is not installed on the target!");
                        Console.WriteLine("[!] Exiting...");
                        return;
                    }
                }
                Console.WriteLine("Unknown arguments!\n" +
                    "Example: SharpRelay.exe uninstall ServiceNameToRemove");
                return;

            }

            if (args.Length != 4)
            {
                Console.WriteLine("Incorrect number of arguments!\n" +
                    "Example: SharpRelay.exe NewServiceName PathToWinDivertDriver OriginalPort DiversionPort");
                return;
            }

            AppDomain.CurrentDomain.UnhandledException += CurrentDomain_UnhandledException;

            Console.WriteLine("[+] Checking if service already exists...");
            if (!ServiceInstaller.ServiceIsInstalled(args[0]))
            {
                Console.WriteLine("[+] Service not found! Creating new service to load driver...");
                try
                {
                    ServiceInstaller.InstallAndStart(args[0], args[0], @"\??\" + args[1]);
                }
                catch
                {
                    Console.WriteLine("[!] Service could not be created!");
                    return;
                }
                
                Console.WriteLine("[+] Service created and started successfully!");
            }
            else
            {
                Console.WriteLine("[!] Service "+args[0]+" already exists... Try again with a different service name!");
                return;
            }


            string filter = "(inbound and tcp.DstPort == " + args[2] + ") or (outbound and tcp.SrcPort == " + args[3] + ")";
            //string filter = "((tcp.DstPort == " + args[2] + " ) or (tcp.SrcPort == " + args[3] + " ))";
            //string filter = "tcp";

            Console.WriteLine("[+] Creating interception handle...");

            var handle = WinDivert.WinDivertOpen(filter, WinDivertLayer.Network, -1000, 0);

            if (handle == IntPtr.Zero || handle == new IntPtr(-1))
            {
                Console.WriteLine("[-] Invalid handle. Failed to open.");
                Console.ReadKey();
                return;
            }

            // Set everything to maximum values.
            WinDivert.WinDivertSetParam(handle, WinDivertParam.QueueLen, 16384);
            WinDivert.WinDivertSetParam(handle, WinDivertParam.QueueTime, 8000);
            WinDivert.WinDivertSetParam(handle, WinDivertParam.QueueSize, 33554432);

            Console.WriteLine("[+] Diverting packets from "+args[2]+" to "+args[3]+"...\n");

            while (s_running)
            {
                SpawnThreads(Environment.ProcessorCount, handle, short.Parse(args[2]), short.Parse(args[3]));
                if (s_running)
                    Console.WriteLine("[i] Threads have apparently crashed... Restarting...");
            }

            WinDivert.WinDivertClose(handle);
        }

        private static void SpawnThreads(int amount, IntPtr handle, short origPort, short divertPort)
        {
            var threads = new List<Thread>();

            for (int i = 0; i < amount; ++i)
            {
                threads.Add(new Thread(() =>
                {
                    RunDiversion(handle, origPort, divertPort);
                }));

                threads.Last().Start();
            }

            foreach (var dt in threads)
            {
                dt.Join();
            }
        }

        private static void CurrentDomain_UnhandledException(object sender, UnhandledExceptionEventArgs e)
        {
            Console.WriteLine($"{(e.ExceptionObject as Exception).Message} at \n{(e.ExceptionObject as Exception).StackTrace}");
        }

        unsafe private static void RunDiversion(IntPtr handle, short origPort, short divertPort)
        {
            IntPtr recvEvent = IntPtr.Zero;
            recvEvent = WinDivertSharp.WinAPI.Kernel32.CreateEvent(IntPtr.Zero, false, false, IntPtr.Zero);

            if (recvEvent == IntPtr.Zero || recvEvent == new IntPtr(-1))
            {
                return;
            }

            var packet = new WinDivertBuffer();

            var addr = new WinDivertAddress();

            uint recvLength = 0;

            NativeOverlapped recvOverlapped;


            uint recvAsyncIoLen = 0;

            ushort[] _v4ReturnPorts = new ushort[ushort.MaxValue + 1];
            Span<byte> payloadBufferPtr = null;

            IPAddress original_client_ip = null;
            IPAddress original_server_ip = null;
            ushort orig_client_tcpSrcPort = 0;

            _originalPort = (ushort)IPAddress.HostToNetworkOrder(origPort);
            _modifiedPort = (ushort)IPAddress.HostToNetworkOrder(divertPort);

            while (s_running)
            {
                try
                {
                    payloadBufferPtr = null;

                    recvLength = 0;
                    addr.Reset();
                    recvAsyncIoLen = 0;

                    recvOverlapped = new NativeOverlapped();

                    recvOverlapped.EventHandle = recvEvent;

                    #region Packet Reading Code

                    if (!WinDivert.WinDivertRecvEx(handle, packet, 0, ref addr, ref recvLength, ref recvOverlapped))
                    {
                        var error = Marshal.GetLastWin32Error();

                        // 997 == ERROR_IO_PENDING
                        if (error != 997)
                        {
                            Console.WriteLine(string.Format("[-] Unknown IO error ID {0} while awaiting overlapped result.", error));
                            Kernel32.CloseHandle(recvEvent);
                            continue;
                        }

                        while (Kernel32.WaitForSingleObject(recvEvent, 1000) == (uint)WaitForSingleObjectResult.WaitTimeout)
                            ;

                        if (!Kernel32.GetOverlappedResult(handle, ref recvOverlapped, ref recvAsyncIoLen, false))
                        {
                            Console.WriteLine("[-] Failed to get overlapped result.");
                            Kernel32.CloseHandle(recvEvent);
                            continue;
                        }

                        recvLength = recvAsyncIoLen;
                    }
                    #endregion Packet Reading Code

                    WinDivertParseResult parseResult = null;
                    parseResult = WinDivert.WinDivertHelperParsePacket(packet, recvLength);

                    if (parseResult == null || parseResult.IPv4Header == null || parseResult.TcpHeader == null)
                        return;

                    if ((parseResult.TcpHeader->Syn == 0x1) && (parseResult.TcpHeader->Ack == 0x0))
                    {
                        original_client_ip = parseResult.IPv4Header->SrcAddr;
                        original_server_ip = parseResult.IPv4Header->DstAddr;
                        orig_client_tcpSrcPort = parseResult.TcpHeader->SrcPort;
                    }

                    if (parseResult.TcpHeader->DstPort == _originalPort)
                    {
                        parseResult.TcpHeader->DstPort = _modifiedPort;
                    }

                    //if (parseResult.TcpHeader->SrcPort == _modifiedPort && parseResult.TcpHeader->DstPort == orig_client_tcpSrcPort)
                    if (parseResult.TcpHeader->SrcPort == _modifiedPort)
                    {
                        parseResult.TcpHeader->SrcPort = _originalPort;
                        //parseResult.TcpHeader->DstPort = orig_client_tcpSrcPort;

                        //parseResult.IPv4Header->SrcAddr = original_server_ip;
                        //parseResult.IPv4Header->DstAddr = original_client_ip;
                    }

                    var sumsCalculated = WinDivert.WinDivertHelperCalcChecksums(packet, recvLength, ref addr, WinDivertChecksumHelperParam.All);

                    WinDivert.WinDivertSendEx(handle, packet, recvLength, 0, ref addr);
                
                }
                catch (Exception ex) { Console.WriteLine($"[!] Exception: {ex.Message} - {ex.StackTrace}"); }
            }
        }

    }
}