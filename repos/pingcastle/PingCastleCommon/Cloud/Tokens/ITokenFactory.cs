//
// Copyright (c) Vincent LE TOUX for Ping Castle. All rights reserved.
// https://www.pingcastle.com
//
// Licensed under the Non-Profit OSL. See LICENSE file in the project root for full license information.
//
using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using PingCastle.Cloud.Credentials;

namespace PingCastle.Cloud.Tokens
{
    using RESTServices.Azure;

    /// <summary>
    /// Interface for handling retrieval and refresh of tokens from authentication providers.
    /// </summary>
    public interface ITokenFactory
    {
        /// <summary>
        /// Retrieves a token from an authentication provider's credential.
        /// </summary>
        /// <typeparam name="T">The type of the Azure service.</typeparam>
        /// <param name="credential">The credential to use for authentication.</param>
        /// <returns>A new <see cref="Token"/>.</returns>
        Task<Token> GetToken<T>(IAzureCredential credential)
            where T : IAzureService;

        /// <summary>
        /// Refreshes the token for a given tenant identity.
        /// </summary>
        /// <typeparam name="T">The type of AzureService to process.</typeparam>
        /// <param name="tenantId">The identity of the tenant.</param>
        /// <param name="token">The <see cref="Token"/> that needs to be refreshed.</param>
        /// <returns>A new <see cref="Token"/>.</returns>
        /// <exception cref="InvalidOperationException">Thrown when the authorization service returns a token with no content.</exception>
        Task<Token> RefreshToken<T>(string tenantId, Token token)
            where T : IAzureService;

        /// <summary>
        /// Retrieves the PRT identities stored on this computer.
        /// </summary>
        /// <returns>A list of <see cref="JwtToken"/> containing the found identities.</returns>
        List<JwtToken> GetRegisteredPRTIdentities();

        /// <summary>
        /// Retrieves an access token from Azure using the authorization code flow.
        /// </summary>
        /// <typeparam name="T">The Azure Service type to work with.</typeparam>
        /// <param name="credential">The credential to retrieve the token for.</param>
        /// <param name="code">The authorization code.</param>
        /// <param name="redirectUri">The redirect URI for the authorization flow.</param>
        /// <param name="codeVerifier">The code verifier.</param>
        /// <returns>A new <see cref="string"/> containing the access token.</returns>
        Task<string> RunGetToken<T>(
            IAzureCredential credential,
            string code,
            string redirectUri,
            string codeVerifier = null)
            where T : IAzureService;
    }
}
