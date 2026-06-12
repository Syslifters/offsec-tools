using System;

namespace PingCastleCommon.Services
{
    using Microsoft.Extensions.Logging;

    /// <summary>
    /// Platform-agnostic logging abstraction.
    /// Platform-specific implementations (Windows event log, file, etc.) are provided by concrete implementations.
    /// </summary>
    public interface IPingCastleLogger : ILogger
    {
        /// <summary>
        /// Write a debug-level log message.
        /// </summary>
        void Debug(string message);

        /// <summary>
        /// Write an information-level log message.
        /// </summary>
        void Information(string message);

        /// <summary>
        /// Write a warning-level log message.
        /// </summary>
        void Warning(string message);

        /// <summary>
        /// Write an error-level log message.
        /// </summary>
        void Error(string message);

        /// <summary>
        /// Write an error-level log message with exception details.
        /// </summary>
        void Error(string message, Exception exception);

        /// <summary>
        /// Write a critical-level log message.
        /// </summary>
        void Critical(string message);

        /// <summary>
        /// Write a critical-level log message with exception details.
        /// </summary>
        void Critical(string message, Exception exception);
    }
}
