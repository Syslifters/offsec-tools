//
// Copyright (c) Ping Castle. All rights reserved.
// https://www.pingcastle.com
//
// Licensed under the Non-Profit OSL. See LICENSE file in the project root for full license information.
//

namespace PingCastle.ADWS
{
    using System.Runtime.CompilerServices;

    public static class ADItemExtensions
    {
        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool IsAccountDisabled(this ADItem item)
        {
            return (item.UserAccountControl & 0x00000002) != 0;
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool IsAccountEnabled(this ADItem item)
        {
            return (item.UserAccountControl & 0x00000002) == 0;
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool IsAccountLockedOut(this ADItem item)
        {
            return (item.UserAccountControl & 0x00000010) != 0;
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool IsPasswordNotRequired(this ADItem item)
        {
            return (item.UserAccountControl & 0x00000020) != 0;
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool IsNormalAccount(this ADItem item)
        {
            return (item.UserAccountControl & 0x00000080) != 0;
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool IsPasswordNeverExpires(this ADItem item)
        {
            return (item.UserAccountControl & 0x00010000) != 0;
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool IsSmartcardRequired(this ADItem item)
        {
            return (item.UserAccountControl & 0x00040000) != 0;
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool IsTrustedForDelegation(this ADItem item)
        {
            return (item.UserAccountControl & 0x00080000) != 0;
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool IsNotDelegated(this ADItem item)
        {
            return (item.UserAccountControl & 0x00100000) != 0;
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool UsesDesKeyOnly(this ADItem item)
        {
            return (item.UserAccountControl & 0x00200000) != 0;
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool DoesNotRequireKerberosPreauthentication(this ADItem item)
        {
            return (item.UserAccountControl & 0x00400000) != 0;
        }

        [MethodImpl(MethodImplOptions.AggressiveInlining)]
        public static bool IsPasswordExpired(this ADItem item)
        {
            return (item.UserAccountControl & 0x01000000) != 0;
        }
    }
}
