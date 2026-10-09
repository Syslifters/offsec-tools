//
// Copyright (c) Ping Castle. All rights reserved.
// https://www.pingcastle.com
//
// Licensed under the Non-Profit OSL. See LICENSE file in the project root for full license information.
//

namespace PingCastle
{
    using Microsoft.Extensions.Options;
    using System;
    using PingCastleCommon.Options;

    public class ADHealthCheckingLicenseSettings
    {
        private static ADHealthCheckingLicenseSettings settings;

        public static ADHealthCheckingLicenseSettings Settings
        {
            get
            {
                if (settings == null)
                {
                    if (ServiceProviderAccessor.IsInitialized)
                    {
                        var options = ServiceProviderAccessor.Current.GetService(typeof(IOptions<LicenseOptions>)) as IOptions<LicenseOptions>;
                        settings = new ADHealthCheckingLicenseSettings();
                        if (options?.Value != null)
                        {
                            settings.License = options.Value.License;
                        }
                    }
                    else
                    {
                        throw new ApplicationException("Could not load license settings.");
                    }
                }

                return settings;
            }
        }

        public string License
        {
            get; set;
        }
    }
}
