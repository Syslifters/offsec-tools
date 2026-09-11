using System.IO;
using PingCastleCommon.Data;

namespace PingCastleCommon.Utility
{
    public interface IHilbertMapGenerator
    {
        bool TryGenerateHilbertImage(Stream stream, NetworkMapDataView view, NetworkMapData data);
    }
}
