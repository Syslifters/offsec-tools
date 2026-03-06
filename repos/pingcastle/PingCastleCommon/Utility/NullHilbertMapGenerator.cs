using System.IO;
using PingCastleCommon.Data;

namespace PingCastleCommon.Utility
{
    public class NullHilbertMapGenerator : IHilbertMapGenerator
    {
        public bool TryGenerateHilbertImage(Stream stream, NetworkMapDataView view, NetworkMapData data)
        {
            return false;
        }
    }
}
