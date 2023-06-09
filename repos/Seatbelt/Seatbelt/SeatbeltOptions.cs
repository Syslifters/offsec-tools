﻿using System.Collections.Generic;

namespace Seatbelt
{
    class SeatbeltOptions
    {
        public SeatbeltOptions(IEnumerable<string> commands, IEnumerable<string> commandGroup, string? outputFile, bool filterResults, bool quietMode, string? computerName, string? userName, string? password)
        {
            Commands = commands;
            CommandGroups = commandGroup;
            OutputFile = outputFile;
            FilterResults = filterResults;
            QuietMode = quietMode;
            ComputerName = computerName;
            UserName = userName;
            Password = password;
        }

        public IEnumerable<string> Commands { get; set; }
        public IEnumerable<string> CommandGroups { get; set; }
        public string? OutputFile { get; set; }
        public bool FilterResults { get; set; }
        public bool QuietMode { get; set; }
        public string? ComputerName { get; set; }
        public string? UserName { get; set; }
        public string? Password { get; set; }
    }
}
