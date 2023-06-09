﻿//
// Copyright (c) Ping Castle. All rights reserved.
// https://www.pingcastle.com
//
// Licensed under the Non-Profit OSL. See LICENSE file in the project root for full license information.
//
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Runtime.InteropServices;
using System.Security.Permissions;
using System.Security.Principal;

namespace PingCastle.RPC
{

    [DebuggerDisplay("{DnsDomainName} {NetbiosDomainName}")]
    public class TrustedDomain
    {
        public string NetbiosDomainName;
        public string DnsDomainName;
        public TrustedDomainFlag Flags;
        public int ParentIndex;
        public int TrustType;
        public int TrustAttributes;
        public SecurityIdentifier DomainSid;
        public Guid DomainGuid;
    }

    [Flags]
    public enum TrustedDomainFlag
    {
        DS_DOMAIN_IN_FOREST = 1,
        DS_DOMAIN_DIRECT_OUTBOUND = 2,
        DS_DOMAIN_TREE_ROOT = 4,
        DS_DOMAIN_PRIMARY = 8,
        DS_DOMAIN_NATIVE_MODE = 16,
        DS_DOMAIN_DIRECT_INBOUND = 32,
    }

    public class nrpc3 : rpcapi
    {

        private static byte[] MIDL_ProcFormatStringx86 = new byte[] {
                0x00,0x48,0x00,0x00,0x00,0x00,0x28,0x00,0x10,0x00,0x31,0x04,0x00,0x00,0x00,0x5c,0x08,0x00,0x08,0x00,0x47,0x04,0x08,0x03,0x01,0x00,0x00,0x00,0x00,
                0x00,0x0b,0x00,0x00,0x00,0x02,0x00,0x48,0x00,0x04,0x00,0x08,0x00,0x13,0x21,0x08,0x00,0xaa,0x00,0x70,0x00,0x0c,0x00,0x08,0x00,0x00
            };

        private static byte[] MIDL_ProcFormatStringx64 = new byte[] {
                0x00,0x48,0x00,0x00,0x00,0x00,0x28,0x00,0x20,0x00,0x31,0x08,0x00,0x00,0x00,0x5c,0x08,0x00,0x08,0x00,0x47,0x04,0x0a,0x03,0x01,0x00,0x00,0x00,0x00,
                0x00,0x00,0x00,0x0b,0x00,0x00,0x00,0x02,0x00,0x48,0x00,0x08,0x00,0x08,0x00,0x13,0x41,0x10,0x00,0x7c,0x00,0x70,0x00,0x18,0x00,0x08,0x00,0x00
        };

        private static byte[] MIDL_TypeFormatStringx86 = new byte[] {
                0x00,0x00,0x12,0x08,0x25,0x5c,0x11,0x04,0xa2,0x00,0x1d,0x00,0x08,0x00,0x01,0x5b,0x15,0x03,0x10,0x00,0x08,0x06,0x06,0x4c,0x00,0xf1,0xff,0x5b,0x1d,
                0x00,0x06,0x00,0x01,0x5b,0x15,0x00,0x06,0x00,0x4c,0x00,0xf4,0xff,0x5c,0x5b,0x1b,0x03,0x04,0x00,0x04,0x00,0xf9,0xff,0x01,0x00,0x08,0x5b,0x17,0x03,
                0x08,0x00,0xf0,0xff,0x02,0x02,0x4c,0x00,0xe0,0xff,0x5c,0x5b,0x16,0x03,0x2c,0x00,0x4b,0x5c,0x46,0x5c,0x00,0x00,0x00,0x00,0x12,0x08,0x25,0x5c,0x46,
                0x5c,0x04,0x00,0x04,0x00,0x12,0x08,0x25,0x5c,0x46,0x5c,0x18,0x00,0x18,0x00,0x12,0x00,0xd0,0xff,0x5b,0x08,0x08,0x08,0x08,0x08,0x08,0x08,0x4c,0x00,
                0x9c,0xff,0x5c,0x5b,0x1b,0x03,0x2c,0x00,0x19,0x00,0x00,0x00,0x01,0x00,0x4b,0x5c,0x48,0x49,0x2c,0x00,0x00,0x00,0x03,0x00,0x00,0x00,0x00,0x00,0x12,
                0x08,0x25,0x5c,0x04,0x00,0x04,0x00,0x12,0x08,0x25,0x5c,0x18,0x00,0x18,0x00,0x12,0x00,0x96,0xff,0x5b,0x4c,0x00,0x9f,0xff,0x5b,0x16,0x03,0x08,0x00,
                0x4b,0x5c,0x46,0x5c,0x04,0x00,0x04,0x00,0x12,0x00,0xc0,0xff,0x5b,0x08,0x08,0x5b,0x00
        };

        private static byte[] MIDL_TypeFormatStringx64 = new byte[] {
                0x00,0x00,0x12,0x08,0x25,0x5c,0x11,0x04,0x74,0x00,0x1d,0x00,0x08,0x00,0x01,0x5b,0x15,0x03,0x10,0x00,0x08,0x06,0x06,0x4c,0x00,0xf1,0xff,0x5b,0x1d,
                0x00,0x06,0x00,0x01,0x5b,0x15,0x00,0x06,0x00,0x4c,0x00,0xf4,0xff,0x5c,0x5b,0x1b,0x03,0x04,0x00,0x04,0x00,0xf9,0xff,0x01,0x00,0x08,0x5b,0x17,0x03,
                0x08,0x00,0xf0,0xff,0x02,0x02,0x4c,0x00,0xe0,0xff,0x5c,0x5b,0x1a,0x03,0x38,0x00,0x00,0x00,0x0e,0x00,0x36,0x36,0x08,0x08,0x08,0x08,0x36,0x4c,0x00,
                0xb9,0xff,0x5b,0x12,0x08,0x25,0x5c,0x12,0x08,0x25,0x5c,0x12,0x00,0xd4,0xff,0x21,0x03,0x00,0x00,0x19,0x00,0x00,0x00,0x01,0x00,0xff,0xff,0xff,0xff,
                0x00,0x00,0x4c,0x00,0xce,0xff,0x5c,0x5b,0x1a,0x03,0x10,0x00,0x00,0x00,0x06,0x00,0x08,0x40,0x36,0x5b,0x12,0x00,0xdc,0xff,0x00
        };

        [StructLayout(LayoutKind.Sequential)]
        internal struct NETLOGON_TRUSTED_DOMAIN_ARRAY
        {
            public int DomainCount;
            public IntPtr Domains;
        }

        [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Unicode)]
        private struct DS_DOMAIN_TRUSTSW
        {
            public IntPtr NetbiosDomainName;
            public IntPtr DnsDomainName;
            public int Flags;
            public int ParentIndex;
            public int TrustType;
            public int TrustAttributes;
            public IntPtr DomainSid;
            public Guid DomainGuid;
        }

        [SecurityPermission(SecurityAction.LinkDemand, Flags = SecurityPermissionFlag.UnmanagedCode)]
        public nrpc3(bool WillUseNullSession = true)
        {
            Guid interfaceId = new Guid(magic(8) + "-" + magic(4) + "-ABCD-EF00-01234567CFFB");
            if (IntPtr.Size == 8)
            {
                InitializeStub(interfaceId, MIDL_ProcFormatStringx64, MIDL_TypeFormatStringx64, "\\pipe\\netlogon");
            }
            else
            {
                InitializeStub(interfaceId, MIDL_ProcFormatStringx86, MIDL_TypeFormatStringx86, "\\pipe\\netlogon");
            }
            UseNullSession = WillUseNullSession;
        }

        [SecurityPermission(SecurityAction.Demand, Flags = SecurityPermissionFlag.UnmanagedCode)]
        ~nrpc3()
        {
            freeStub();
        }

        [SecurityPermission(SecurityAction.Demand, Flags = SecurityPermissionFlag.UnmanagedCode)]
        public Int32 DsrEnumerateDomainTrusts(string server, int flag, out List<TrustedDomain> domains)
        {
            IntPtr result = IntPtr.Zero;
            domains = null;
            IntPtr intptrServer = Marshal.StringToHGlobalUni(server);
            NETLOGON_TRUSTED_DOMAIN_ARRAY output = new NETLOGON_TRUSTED_DOMAIN_ARRAY();
            try
            {
                if (IntPtr.Size == 8)
                {
                    result = NativeMethods.NdrClientCall2x64(GetStubHandle(), GetProcStringHandle(0), intptrServer, flag, ref output);
                }
                else
                {
                    GCHandle handle = GCHandle.Alloc(output, GCHandleType.Pinned);
                    IntPtr tempValuePointer = handle.AddrOfPinnedObject();
                    try
                    {
                        result = CallNdrClientCall2x86(0, intptrServer, new IntPtr((int)flag), tempValuePointer);
                        // each pinvoke work on a copy of the arguments (without an out specifier)
                        // get back the data
                        output = (NETLOGON_TRUSTED_DOMAIN_ARRAY)Marshal.PtrToStructure(tempValuePointer, typeof(NETLOGON_TRUSTED_DOMAIN_ARRAY));
                    }
                    finally
                    {
                        handle.Free();
                    }
                }
            }
            catch (SEHException)
            {
                return Marshal.GetExceptionCode();
            }
            finally
            {
                if (intptrServer != IntPtr.Zero)
                    Marshal.FreeHGlobal(intptrServer);
            }
            domains = DomainArrayToTrustedDomainList(output);
            return (int)result.ToInt64();
        }

        [SecurityPermission(SecurityAction.Demand, Flags = SecurityPermissionFlag.UnmanagedCode)]
        private List<TrustedDomain> DomainArrayToTrustedDomainList(NETLOGON_TRUSTED_DOMAIN_ARRAY trustedDomainArray)
        {
            List<TrustedDomain> output = new List<TrustedDomain>();
            int size = Marshal.SizeOf(typeof(DS_DOMAIN_TRUSTSW));
            for (int i = 0; i < trustedDomainArray.DomainCount; i++)
            {
                DS_DOMAIN_TRUSTSW trust = (DS_DOMAIN_TRUSTSW)Marshal.PtrToStructure(new IntPtr(trustedDomainArray.Domains.ToInt64() + size * i), typeof(DS_DOMAIN_TRUSTSW));
                TrustedDomain domain = new TrustedDomain();
                if (trust.DnsDomainName != IntPtr.Zero)
                {
                    domain.DnsDomainName = Marshal.PtrToStringUni(trust.DnsDomainName);
                    FreeMemory(trust.DnsDomainName);
                }
                if (trust.NetbiosDomainName != IntPtr.Zero)
                {
                    domain.NetbiosDomainName = Marshal.PtrToStringUni(trust.NetbiosDomainName);
                    FreeMemory(trust.NetbiosDomainName);
                }
                domain.Flags = (TrustedDomainFlag)trust.Flags;
                domain.ParentIndex = trust.ParentIndex;
                domain.TrustAttributes = trust.TrustAttributes;
                domain.TrustType = trust.TrustType;
                domain.DomainGuid = trust.DomainGuid;
                if (trust.DomainSid != IntPtr.Zero)
                {
                    domain.DomainSid = new SecurityIdentifier(trust.DomainSid);
                    FreeMemory(trust.DomainSid);
                }
                output.Add(domain);
            }
            FreeMemory(trustedDomainArray.Domains);
            return output;
        }

    }
}
