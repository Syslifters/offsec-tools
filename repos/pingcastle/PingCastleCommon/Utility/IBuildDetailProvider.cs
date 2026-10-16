using System;
using System.Reflection;

namespace PingCastleCommon.Utility
{
    public interface IBuildDetailProvider
    {
        DateTime GetBuildDateTime(Assembly assembly);
        DateTime GetBuildDateTime(string path);
    }
}
