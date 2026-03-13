//
// Copyright (c) Vincent LE TOUX for Ping Castle. All rights reserved.
// https://www.pingcastle.com
//
// Licensed under the Non-Profit OSL. See LICENSE file in the project root for full license information.
//
namespace PingCastle.Cloud.Credentials
{
    using PingCastle.Cloud.Tokens;

    public class UserCredential : CredentialBase
    {
        public UserCredential(ITokenFactory tokenFactory) : base(tokenFactory)
        {
        }

        public UserCredential(string tenantid, ITokenFactory tokenFactory) : base(tenantid, tokenFactory)
        {
        }
    }
}
