const osList = [];

function editDistance(s1, s2) {
  s1 = s1.toLowerCase();
  s2 = s2.toLowerCase();

  var costs = new Array();
  for (var i = 0; i <= s1.length; i++) {
    var lastValue = i;
    for (var j = 0; j <= s2.length; j++) {
      if (i == 0)
        costs[j] = j;
      else {
        if (j > 0) {
          var newValue = costs[j - 1];
          if (s1.charAt(i - 1) != s2.charAt(j - 1))
            newValue = Math.min(Math.min(newValue, lastValue),
              costs[j]) + 1;
          costs[j - 1] = lastValue;
          lastValue = newValue;
        }
      }
    }
    if (i > 0)
      costs[s2.length] = lastValue;
  }
  return costs[s2.length];
}

function similarity(s1, s2) {
  var longer = s1;
  var shorter = s2;
  if (s1.length < s2.length) {
    longer = s2;
    shorter = s1;
  }
  var longerLength = longer.length;
  if (longerLength == 0) {
    return 1.0;
  }
  return (longerLength - editDistance(longer, shorter)) / parseFloat(longerLength);
}

export default {
    browser: (client) => {
        const brand = client.info.browser_name.toLowerCase()
        return `https://raw.githubusercontent.com/alrra/browser-logos/main/src/${brand}/${brand}_128x128.png`
    },

    os: (client) => {
        const brand = client.info.os_name.toLowerCase()
        const match = osList.reduce((prev, curr) => {
            const sim = similarity(brand, curr.slug.toLowerCase())
            const dis = editDistance(brand, curr.slug.toLowerCase())
            console.log(brand, curr, sim, dis)
            if (sim > prev.sim || (sim === prev.sim && dis < prev.dis)) {
              return {sim, dis, curr}
            } else if (sim == prev.sim && dis == prev.dis) {
              if (brand.includes(curr.slug.toLowerCase()) || curr.slug.toLowerCase().includes(brand)) {
                return {sim, dis, curr}
              }
            }
            return prev
        }, {sim: 0, dis: 999, curr: {}})
        return `https://raw.githubusercontent.com/EgoistDeveloper/operating-system-logos/master/src/128x128/${match.curr.code}.png`
    }
}

osList.push(...[
  {
    "code": "AIX",
    "name": "AIX",
    "slug": "aix"
  },
  {
    "code": "AND",
    "name": "Android",
    "slug": "android"
  },
  {
    "code": "AMG",
    "name": "AmigaOS",
    "slug": "amigaos"
  },
  {
    "code": "ATV",
    "name": "tvOS",
    "slug": "tvos"
  },
  {
    "code": "ARL",
    "name": "Arch Linux",
    "slug": "arch-linux"
  },
  {
    "code": "BTR",
    "name": "BackTrack",
    "slug": "backtrack"
  },
  {
    "code": "SBA",
    "name": "Bada",
    "slug": "bada"
  },
  {
    "code": "BEO",
    "name": "BeOS",
    "slug": "beos"
  },
  {
    "code": "BLB",
    "name": "BlackBerry OS",
    "slug": "blackberry-os"
  },
  {
    "code": "QNX",
    "name": "BlackBerry Tablet OS",
    "slug": "blackberry-tablet-os"
  },
  {
    "code": "CAI",
    "name": "Caixa Mágica",
    "slug": "caixa-mgica"
  },
  {
    "code": "CES",
    "name": "CentOS",
    "slug": "centos"
  },
  {
    "code": "COS",
    "name": "Chrome OS",
    "slug": "chrome-os"
  },
  {
    "code": "CYN",
    "name": "CyanogenMod",
    "slug": "cyanogenmod"
  },
  {
    "code": "DEB",
    "name": "Debian",
    "slug": "debian"
  },
  {
    "code": "DEE",
    "name": "Deepin",
    "slug": "deepin"
  },
  {
    "code": "DFB",
    "name": "DragonFly",
    "slug": "dragonfly"
  },
  {
    "code": "FED",
    "name": "Fedora",
    "slug": "fedora"
  },
  {
    "code": "FOS",
    "name": "Firefox OS",
    "slug": "firefox-os"
  },
  {
    "code": "FIR",
    "name": "Fire OS",
    "slug": "fire-os"
  },
  {
    "code": "BSD",
    "name": "FreeBSD",
    "slug": "freebsd"
  },
  {
    "code": "FYD",
    "name": "FydeOS",
    "slug": "fydeos"
  },
  {
    "code": "GNT",
    "name": "Gentoo",
    "slug": "gentoo"
  },
  {
    "code": "GTV",
    "name": "Google TV",
    "slug": "google-tv"
  },
  {
    "code": "HPX",
    "name": "HP-UX",
    "slug": "hp-ux"
  },
  {
    "code": "HAI",
    "name": "Haiku OS",
    "slug": "haiku-os"
  },
  {
    "code": "IPA",
    "name": "iPadOS",
    "slug": "ipados"
  },
  {
    "code": "HAR",
    "name": "HarmonyOS",
    "slug": "harmonyos"
  },
  {
    "code": "KOS",
    "name": "KaiOS",
    "slug": "kaios"
  },
  {
    "code": "KNO",
    "name": "Knoppix",
    "slug": "knoppix"
  },
  {
    "code": "KBT",
    "name": "Kubuntu",
    "slug": "kubuntu"
  },
  {
    "code": "LIN",
    "name": "GNU/Linux",
    "slug": "gnulinux"
  },
  {
    "code": "LBT",
    "name": "Lubuntu",
    "slug": "lubuntu"
  },
  {
    "code": "MAC",
    "name": "Mac",
    "slug": "mac"
  },
  {
    "code": "MAE",
    "name": "Maemo",
    "slug": "maemo"
  },
  {
    "code": "MAG",
    "name": "Mageia",
    "slug": "mageia"
  },
  {
    "code": "MDR",
    "name": "Mandriva",
    "slug": "mandriva"
  },
  {
    "code": "SMG",
    "name": "MeeGo",
    "slug": "meego"
  },
  {
    "code": "MIN",
    "name": "Mint",
    "slug": "mint"
  },
  {
    "code": "MOR",
    "name": "MorphOS",
    "slug": "morphos"
  },
  {
    "code": "NBS",
    "name": "NetBSD",
    "slug": "netbsd"
  },
  {
    "code": "WII",
    "name": "Nintendo",
    "slug": "nintendo"
  },
  {
    "code": "NDS",
    "name": "Nintendo Mobile",
    "slug": "nintendo-mobile"
  },
  {
    "code": "OS2",
    "name": "OS/2",
    "slug": "os2"
  },
  {
    "code": "OBS",
    "name": "OpenBSD",
    "slug": "openbsd"
  },
  {
    "code": "OWR",
    "name": "OpenWrt",
    "slug": "openwrt"
  },
  {
    "code": "PCL",
    "name": "PCLinuxOS",
    "slug": "pclinuxos"
  },
  {
    "code": "PSP",
    "name": "PlayStation Portable",
    "slug": "playstation-portable"
  },
  {
    "code": "PS3",
    "name": "PlayStation",
    "slug": "playstation"
  },
  {
    "code": "RAS",
    "name": "Raspberry Pi",
    "slug": "raspberrypi"
  },
  {
    "code": "RHT",
    "name": "Red Hat",
    "slug": "red-hat"
  },
  {
    "code": "ROS",
    "name": "RISC OS",
    "slug": "risc-os"
  },
  {
    "code": "ROK",
    "name": "Roku OS",
    "slug": "roku-os"
  },
  {
    "code": "RSO",
    "name": "Rosa",
    "slug": "rosa"
  },
  {
    "code": "REM",
    "name": "Remix OS",
    "slug": "remix-os"
  },
  {
    "code": "REX",
    "name": "REX",
    "slug": "rex"
  },
  {
    "code": "SAB",
    "name": "Sabayon",
    "slug": "sabayon"
  },
  {
    "code": "SSE",
    "name": "SUSE",
    "slug": "suse"
  },
  {
    "code": "SAF",
    "name": "Sailfish OS",
    "slug": "sailfish-os"
  },
  {
    "code": "SLW",
    "name": "Slackware",
    "slug": "slackware"
  },
  {
    "code": "SOS",
    "name": "Solaris",
    "slug": "solaris"
  },
  {
    "code": "SYL",
    "name": "Syllable",
    "slug": "syllable"
  },
  {
    "code": "SYM",
    "name": "Symbian",
    "slug": "symbian"
  },
  {
    "code": "SYS",
    "name": "Symbian OS",
    "slug": "symbian-os"
  },
  {
    "code": "S40",
    "name": "Symbian OS Series 40",
    "slug": "symbian-os-series-40"
  },
  {
    "code": "S60",
    "name": "Symbian OS Series 60",
    "slug": "symbian-os-series-60"
  },
  {
    "code": "SY3",
    "name": "Symbian^3",
    "slug": "symbian3"
  },
  {
    "code": "TDX",
    "name": "ThreadX",
    "slug": "threadx"
  },
  {
    "code": "TIZ",
    "name": "Tizen",
    "slug": "tizen"
  },
  {
    "code": "UBT",
    "name": "Ubuntu",
    "slug": "ubuntu"
  },
  {
    "code": "WAS",
    "name": "watchOS",
    "slug": "watchos"
  },
  {
    "code": "WHS",
    "name": "Whale OS",
    "slug": "whale-os"
  },
  {
    "code": "WIN",
    "name": "Windows",
    "slug": "windows"
  },
  {
    "code": "WCE",
    "name": "Windows CE",
    "slug": "windows-ce"
  },
  {
    "code": "WIO",
    "name": "Windows IoT",
    "slug": "windows-iot"
  },
  {
    "code": "WMO",
    "name": "Windows Mobile",
    "slug": "windows-mobile"
  },
  {
    "code": "WPH",
    "name": "Windows Phone",
    "slug": "windows-phone"
  },
  {
    "code": "WRT",
    "name": "Windows RT",
    "slug": "windows-rt"
  },
  {
    "code": "XBX",
    "name": "Xbox",
    "slug": "xbox"
  },
  {
    "code": "XBT",
    "name": "Xubuntu",
    "slug": "xubuntu"
  },
  {
    "code": "YNS",
    "name": "YunOs",
    "slug": "yunos"
  },
  {
    "code": "IOS",
    "name": "iOS",
    "slug": "ios"
  },
  {
    "code": "POS",
    "name": "palmOS",
    "slug": "palmos"
  },
  {
    "code": "WOS",
    "name": "webOS",
    "slug": "webos"
  }
])