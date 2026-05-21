#nullable enable
namespace PingCastleCommon;

using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using PingCastle.Scanners;
using PingCastleCommon.Options;
using PingCastleCommon.Services;
using System;
using System.Linq;

/// <summary>
/// Extension methods for configuring PingCastleCommon services and options.
/// This class provides platform-agnostic configuration that can be used by
/// any consumer of PingCastleCommon (console apps, web apps, etc.)
/// </summary>
public static class ConfigurationExtensions
{
    /// <summary>
    /// Registers all PingCastleCommon configuration options from the provided IConfiguration.
    /// This should be called after the base AddPingCastleCommonServices() to configure
    /// the options that common services depend on.
    /// </summary>
    /// <param name="services">The service collection to register options in</param>
    /// <param name="configuration">The application configuration</param>
    /// <returns>The service collection for chaining</returns>
    public static IServiceCollection AddPingCastleCommonOptions(
        this IServiceCollection services,
        IConfiguration configuration)
    {
        // Register all options from configuration sections
        services.Configure<LicenseOptions>(configuration.GetSection(LicenseOptions.SectionName));
        services.Configure<EncryptionOptions>(configuration.GetSection(EncryptionOptions.SectionName));
        services.Configure<HoneyPotOptions>(configuration.GetSection(HoneyPotOptions.SectionName));
        services.Configure<InfrastructureOptions>(configuration.GetSection(InfrastructureOptions.SectionName));
        services.Configure<CustomRulesOptions>(configuration.GetSection(CustomRulesOptions.SectionName));
        services.Configure<SmtpOptions>(configuration.GetSection(SmtpOptions.SectionName));
        services.Configure<BrandOptions>(configuration.GetSection(BrandOptions.SectionName));

        return services;
    }

    /// <summary>
    /// Registers PingCastleCommon services including healthcheck rules.
    /// Rules require IIdentityProvider to be registered separately by the caller.
    /// </summary>
    /// <param name="services">The service collection to register services in</param>
    /// <returns>The service collection for chaining</returns>
    public static IServiceCollection AddPingCastleCommonServices(
        this IServiceCollection services)
    {
        // Register scanners that require dependency injection.
        services.AddTransient<KerberosChecksumVulnerabilityScanner>();
        services.AddTransient<SmbHotFixVulnerabilityScanner>();

        // Resource management
        services.AddSingleton<IResourceManagerProvider, ResourceManagerProvider>();

        // Register rules that require dependency injection
        var cloudRuleBaseType = typeof(PingCastle.Rules.RuleBase<PingCastle.Cloud.Data.HealthCheckCloudData>);
        var healthcheckRuleBaseType = typeof(PingCastle.Rules.RuleBase<PingCastle.Healthcheck.HealthcheckData>);

        foreach (var type in AppDomain.CurrentDomain.GetAssemblies()
            .SelectMany(a =>
            {
                try
                {
                    return a.GetExportedTypes();
                }
                catch
                {
                    return [];
                }
            })
            .Where(t => !t.IsAbstract && t.IsPublic &&
                (cloudRuleBaseType.IsAssignableFrom(t) || healthcheckRuleBaseType.IsAssignableFrom(t))))
        {
            services.AddTransient(type);
        }

        return services;
    }

    /// <summary>
    /// Convenience methods to get option values from the service provider.
    /// These allow options to be accessed without directly using IOptions{T}.
    /// </summary>

    public static LicenseOptions? GetLicenseSettings(this IServiceProvider serviceProvider)
    {
        var options = serviceProvider.GetService<Microsoft.Extensions.Options.IOptions<LicenseOptions>>();
        return options?.Value;
    }

    public static EncryptionOptions? GetEncryptionSettings(this IServiceProvider serviceProvider)
    {
        var options = serviceProvider.GetService<Microsoft.Extensions.Options.IOptions<EncryptionOptions>>();
        return options?.Value;
    }

    public static HoneyPotOptions? GetHoneyPotSettings(this IServiceProvider serviceProvider)
    {
        var options = serviceProvider.GetService<Microsoft.Extensions.Options.IOptions<HoneyPotOptions>>();
        return options?.Value;
    }

    public static InfrastructureOptions? GetInfrastructureSettings(this IServiceProvider serviceProvider)
    {
        var options = serviceProvider.GetService<Microsoft.Extensions.Options.IOptions<InfrastructureOptions>>();
        return options?.Value;
    }

    public static CustomRulesOptions? GetCustomRulesSettings(this IServiceProvider serviceProvider)
    {
        var options = serviceProvider.GetService<Microsoft.Extensions.Options.IOptions<CustomRulesOptions>>();
        return options?.Value;
    }

    public static SmtpOptions? GetSmtpSettings(this IServiceProvider serviceProvider)
    {
        var options = serviceProvider.GetService<Microsoft.Extensions.Options.IOptions<SmtpOptions>>();
        return options?.Value;
    }

    public static BrandOptions? GetBrandSettings(this IServiceProvider serviceProvider)
    {
        var options = serviceProvider.GetService<Microsoft.Extensions.Options.IOptions<BrandOptions>>();
        return options?.Value;
    }
}
