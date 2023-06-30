﻿using SnaffCore.Concurrency;
using System;
using System.Collections.Generic;
using System.Text.RegularExpressions;
using static SnaffCore.Config.Options;

namespace SnaffCore.Classifiers
{
    public class TextClassifier
    {
        private ClassifierRule ClassifierRule { get; set; }
        public TextClassifier(ClassifierRule inRule)
        {
            this.ClassifierRule = inRule;
        }

        private BlockingMq Mq { get; set; } = BlockingMq.GetMq();

        // Methods for classification
        internal TextResult TextMatch(string input)
        {
            foreach (Regex regex in ClassifierRule.Regexes)
            {
                try
                {
                    if (regex.IsMatch(input))
                    {
                        return new TextResult()
                        {
                            MatchedStrings = new List<string>() { regex.ToString() },
                            MatchContext = GetContext(input, regex)
                        };
                    }
                }
                catch (Exception e)
                {
                    Mq.Error(e.ToString());
                }
            }
            return null;
        }
        internal string GetContext(string original, string matchString)
        {
            try
            {
                int contextBytes = MyOptions.MatchContextBytes;
                if (contextBytes == 0)
                {
                    return "";
                }

                if (original.Length <= (contextBytes * 2))
                {
                    return original;
                }

                int foundIndex = original.IndexOf(matchString, StringComparison.OrdinalIgnoreCase);

                int contextStart = SubtractWithFloor(foundIndex, contextBytes, 0);
                string matchContext = "";
                if (contextBytes > 0) matchContext = original.Substring(contextStart, contextBytes * 2);

                return Regex.Escape(matchContext);
            }
            catch (ArgumentOutOfRangeException)
            {
                return original;
            }
            catch (Exception e)
            {
                Mq.Trace(e.ToString());
                return "";
            }
        }
        internal string GetContext(string original, Regex matchRegex)
        {
            try
            {
                int contextBytes = MyOptions.MatchContextBytes;
                if (contextBytes == 0)
                {
                    return "";
                }

                if ((original.Length < 6) || (original.Length < contextBytes * 2))
                {
                    return original;
                }

                int foundIndex = matchRegex.Match(original).Index;

                int contextStart = SubtractWithFloor(foundIndex, contextBytes, 0);
                string matchContext = "";

                if (original.Length <= (contextStart + (contextBytes * 2)))
                {
                    return Regex.Escape(original.Substring(contextStart));
                }

                if (contextBytes > 0) matchContext = original.Substring(contextStart, contextBytes * 2);

                return Regex.Escape(matchContext);
            }
            catch (Exception e)
            {
                this.Mq.Error(e.ToString());
            }

            return "";
        }

        internal int SubtractWithFloor(int num1, int num2, int floor)
        {
            int result = num1 - num2;
            if (result <= floor) return floor;
            return result;
        }
    }

    public class TextResult
    {
        public List<string> MatchedStrings { get; set; }
        public string MatchContext { get; set; }
    }
}