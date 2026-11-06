namespace PingCastle.ADWS;

using System;
using System.Collections.Generic;
using System.IO;
using System.Security.AccessControl;

public interface IFileConnection : IDisposable
{
    bool IsDirectory(string path);
    DirectorySecurity GetDirectorySecurity(string path);
    FileSecurity GetFileSecurity(string path);
    IEnumerable<string> GetSubDirectories(string path);
    bool FileExists(string path);
    bool DirectoryExists(string path);
    string GetShortName(string path);
    Stream GetFileStream(string path);
    DateTime GetLastWriteTime(string path);
    string PathCombine(string path1, string path2);
    List<string> GetAllSubDirectories(string path);
    List<string> GetAllSubFiles(string path);

    void ThreadInitialization();

    /// <summary>
    /// Provides impersonation for file access. Wrap file functionality in a call to this method.
    /// </summary>
    /// <typeparam name="T">The functional return type.</typeparam>
    /// <param name="func">The function to execute inside the impersonation call.</param>
    /// <returns>The value of the function called.</returns>
    T RunImpersonatedIfNeeded<T>(Func<T> func);

    /// <summary>
    /// Provides impersonation for file access. Wrap file functionality in a call to this method.
    /// </summary>
    /// <param name="action">The action to execute with impersonation.</param>
    void RunImpersonatedIfNeeded(Action action);
}