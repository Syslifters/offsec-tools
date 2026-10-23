//
// Copyright (c) Vincent LE TOUX for Ping Castle. All rights reserved.
// https://www.pingcastle.com
//
// Licensed under the Non-Profit OSL. See LICENSE file in the project root for full license information.
//

namespace PingCastle.Cloud.Credentials;

using PingCastle.Cloud.RESTServices.Azure;
using PingCastle.Cloud.Tokens;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

public abstract class CredentialBase : IAzureCredential
{
    private readonly Dictionary<Type, Token> _cache = new Dictionary<Type, Token>();
    private readonly ITokenFactory _tokenFactory;

    public CredentialBase(ITokenFactory tokenFactory)
        : this(null, tokenFactory)
    {
    }

    public CredentialBase(string tenantId, ITokenFactory tokenFactory)
    {
        _tenantId = tenantId;
        _tokenFactory = tokenFactory;
    }

    public Token LastTokenQueried { get; protected set; }

    public bool ForceRefreshByRefreshToken { get; set; }

    public async Task<Token> GetToken<T>() where T : IAzureService
    {
        if (_cache.ContainsKey(typeof(T)))
        {
            var cachedToken = _cache[typeof(T)];

            const int networkLatency = 5;
            var expiresOn = DateTimeOffset.FromUnixTimeSeconds(cachedToken.expires_on).AddSeconds(-networkLatency);

            if (expiresOn <= DateTime.UtcNow || ForceRefreshByRefreshToken)
            {
                cachedToken = await _tokenFactory.RefreshToken<T>(_tenantId, cachedToken);
                UpdateTokenCache<T>(cachedToken);
            }

            return cachedToken;
        }

        var newToken = await _tokenFactory.GetToken<T>(this);
        UpdateTokenCache<T>(newToken);

        return newToken;
    }

    private void UpdateTokenCache<T>(Token token) where T : IAzureService
    {
        if (token.expires_on == 0)
        {
            token.expires_on = (uint)((DateTimeOffset)DateTime.UtcNow.AddSeconds(token.expires_in)).ToUnixTimeSeconds();
        }

        LastTokenQueried = token;
        _cache[typeof(T)] = token;
    }

    string _tenantId;

    public string Tenantid
    {
        get
        {
            if (string.IsNullOrEmpty(_tenantId) && _cache.Count > 0)
            {
                var token = _cache.Values.FirstOrDefault();
                _tenantId = token.ToJwtToken().tid;
            }

            return _tenantId;
        }
    }

    string _tenantidToQuery;
    public string TenantidToQuery
    {
        get => _tenantidToQuery;
        set
        {
            if (_tenantidToQuery == value)
            {
                return;
            }

            if (string.IsNullOrEmpty(value))
            {
                return;
            }

            _tenantidToQuery = value;
            _cache.Clear();
        }
    }
}