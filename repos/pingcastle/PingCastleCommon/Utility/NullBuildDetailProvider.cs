using System;
using System.IO;
using System.Reflection;

namespace PingCastleCommon.Utility
{
    public class NullBuildDetailProvider : IBuildDetailProvider
    {
        public DateTime GetBuildDateTime(Assembly assembly)
        {
            var path = assembly.Location;
            return GetBuildDateTime(path);
        }

        public DateTime GetBuildDateTime(string path)
        {
            if (File.Exists(path))
            {
                return File.GetLastWriteTimeUtc(path);
            }

            return DateTime.MinValue;
        }
    }
}
