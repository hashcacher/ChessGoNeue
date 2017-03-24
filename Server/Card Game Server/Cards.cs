using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Collections;

namespace CardGameServer
{
    class Cards
    {
        public static object[] cardList =
        {
            // Name - Description - Author - Cost red - Cost green - Cost blue - Cost yellow - Cost neutral - Attack - HP - Effects
            new object[] { "Little Robot", "I'm a little robot, short and stout", "Sonia", 0, 0, 0, 2, 1, 2000, 1000, 0 },
            new object[] { "Bird Rabbit", "I'm shooting a bow!", "Sonia", 0, 0, 0, 2, 0, 50, 50, 0 },
            new object[] { "Owl Shaman", "Come get sacrificed", "Sonia", 0, 0, 0, 2, 0, 999, 999, 0 },
            new object[] { "Dragon", "Rawr", "Sonia", 0, 0, 0, 2, 0, 999, 999, 0 },
            new object[] { "Griffin", "What is he up to?", "Sonia", 0, 0, 0, 2, 0, 999, 999, 0 },
            new object[] { "Puppy eyes", "Aww what a cute little animal... BRAAWR", "Sonia", 0, 0, 0, 2, 0, 999, 999, 0 },
            new object[] { "Bombs Away!", "A flock of diving birds", "Sonia", 0, 0, 0, 2, 0, 999, 999, 0 },
            new object[] { "Spare parts", "Spare parts from the Ikea shelf unit", "Sonia", 0, 0, 0, 2, 0, 999, 999, 0 },
            new object[] { "Call to the pack", "Here comes danger", "Sonia", 0, 0, 0, 2, 0, 999, 999, 0 },
            new object[] { "Giant tree monster", "You humans fucked up", "Sonia", 0, 0, 0, 2, 0, 999, 999, 0 },
            new object[] { "Racoon Thief", "He ain't no Robin Hood", "Sonia", 0, 0, 0, 2, 0, 999, 999, 0 }
        };
    }
}
