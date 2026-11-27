namespace PingCastleCommon.Options
{
    public class SmtpOptions
    {
        public const string SectionName = "Smtp";

        public string From { get; set; } = string.Empty;

        public string DeliveryMethod { get; set; } = "Network";

        public string Host { get; set; } = string.Empty;

        public int Port { get; set; } = 25;

        public string UserName { get; set; } = string.Empty;

        public string Password { get; set; } = string.Empty;
    }
}
