using System;
using System.Runtime.InteropServices;

namespace PingCastle.RPC
{
    [StructLayout(LayoutKind.Sequential)]
    public struct SecHandle //=PCtxtHandle
    {
        IntPtr dwLower; // ULONG_PTR translates to IntPtr not to uint
        IntPtr dwUpper; // this is crucial for 64-Bit Platforms
    }
}
