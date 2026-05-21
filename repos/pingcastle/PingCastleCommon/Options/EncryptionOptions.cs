using System.Collections.Generic;

namespace PingCastleCommon.Options
{
    public class EncryptionOptions
    {
        public const string SectionName = "Encryption";

        public string EncryptionKey { get; set; } = string.Empty;

        public List<KeySettingsOption> RSAKeys { get; set; } = new();
    }

    public class KeySettingsOption
    {
        public string Name { get; set; } = string.Empty;

        public string? PublicKey { get; set; }

        public string? PrivateKey { get; set; }
    }
}
