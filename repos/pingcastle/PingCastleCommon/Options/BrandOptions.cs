namespace PingCastleCommon.Options
{
    /// <summary>
    /// Configuration options for brand customization in reports.
    /// These settings allow customers with a license to customize the visual appearance
    /// and scripting behavior of generated reports.
    /// </summary>
    public class BrandOptions
    {
        /// <summary>
        /// Configuration section name in appsettings.json
        /// </summary>
        public const string SectionName = "Brand";

        /// <summary>
        /// Gets or sets the base64-encoded brand logo to be displayed in reports.
        /// </summary>
        public string BrandLogo { get; set; }

        /// <summary>
        /// Gets or sets custom CSS styles to be injected into reports.
        /// </summary>
        public string BrandCss { get; set; }

        /// <summary>
        /// Gets or sets custom JavaScript code to be injected into reports.
        /// </summary>
        public string BrandJs { get; set; }
    }
}
