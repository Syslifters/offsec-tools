namespace PingCastle.Rules
{
    public interface IMakeAnalysisVisible<T>
    {
        int? RunAnalysis(T healthcheckData);
    }
}