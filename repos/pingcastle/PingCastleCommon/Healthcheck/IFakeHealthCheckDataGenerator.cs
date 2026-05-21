namespace PingCastleCommon.Healthcheck;

using PingCastle.Data;
using PingCastle.Healthcheck;

/// <summary>
/// Interface for generating fake health check data for testing and demonstration purposes.
/// </summary>
public interface IFakeHealthCheckDataGenerator
{
    PingCastleReportCollection<HealthcheckData> GenerateData();
    PingCastleReportCollection<HealthcheckData> GenerateData(FakeHealthCheckDataGeneratorModel model);
    PingCastleReportCollection<HealthcheckData> GenerateForest(int maxDomain = 15);
    HealthcheckData GenerateSingleReport(FakeHealthCheckDataGeneratorDomainModel model);
}
