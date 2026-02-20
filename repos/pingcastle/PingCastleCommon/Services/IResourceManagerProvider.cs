namespace PingCastleCommon.Services;

using System.Resources;

public interface IResourceManagerProvider
{
    ResourceManager GetHealthCheckRuleResourceManager();
}