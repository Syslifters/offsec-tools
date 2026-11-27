// Copyright (c) Ping Castle. All rights reserved.
// https://www.pingcastle.com
//
// Licensed under the Non-Profit OSL. See LICENSE file in the project root for full license information.
//

namespace PingCastle.Rules
{
    public class RuleAnalysisDecorator<T>
    {
        private readonly RuleBase<T> _rule;

        public RuleAnalysisDecorator(RuleBase<T> rule)
        {
            _rule = rule;
        }

        public int? RunAnalysis(T healthcheckData)
        {
            return _rule.AnalyzeDataNew(healthcheckData);
        }
    }
}
