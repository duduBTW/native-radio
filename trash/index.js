const path = require("path");
const Realm = require("realm");

const args = {};

process.argv.slice(2).forEach((arg) => {
  const [key, value] = arg.replace(/^--/, "").split("=");
  args[key] = value;
});

async function parseLazerDatabase() {
  const databasePath = args.lazer;
  if (!databasePath) {
    return;
  }
  const currentDir = databasePath.replaceAll("\\", "/");
  const realm = await Realm.open({
    path: currentDir + "/client.realm",
    readOnly: true,
    schemaVersion: 23,
  });
  const beatmapSets = realm.objects("BeatmapSet");

  const songTable = new Map();
  const audioTable = new Map();

  for (const beatmapSet of beatmapSets) {
    try {
      const beatmaps = beatmapSet.Beatmaps;

      for (const beatmap of beatmaps) {
        try {
          const song = {
            audio: "",
            osuFile: "",
            path: "",
            ctime: "",
            dateAdded: beatmapSet.DateAdded,
            title: beatmap.Metadata.Title,
            artist: beatmap.Metadata.Artist,
            creator: beatmap.Metadata.Author?.Username ?? "No Creator",
            bpm: [[beatmap.BPM]],
            duration: beatmap.Length,
            diffs: [beatmap.DifficultyName ?? "Unknown difficulty"],
          };

          song.osuFile = path.join(
            currentDir,
            "files",
            beatmap.Hash[0],
            beatmap.Hash.substring(0, 2),
            beatmap.Hash
          );

          const songHash = beatmapSet.Files.find(
            (file) =>
              file.Filename.toLowerCase() ===
              beatmap.Metadata.AudioFile.toLowerCase()
          )?.File.Hash;

          if (songHash) {
            song.audio = path.join(
              currentDir,
              "files",
              songHash[0],
              songHash.substring(0, 2),
              songHash
            );
          }

          const existingSong = songTable.get(song.audio);
          if (existingSong) {
            existingSong.diffs.push(song.diffs[0]);
            continue;
          }

          song.path = song.audio;

          const bgHash = beatmapSet.Files.find(
            (file) => file.Filename === beatmap.Metadata.BackgroundFile
          )?.File.Hash;

          if (bgHash) {
            song.bg = path.join(
              currentDir,
              "files",
              bgHash[0],
              bgHash.substring(0, 2),
              bgHash
            );
          }

          song.beatmapSetID = beatmapSet.OnlineID;

          songTable.set(song.audio, song);
          audioTable.set(song.audio, {
            songID: song.audio,
            path: song.audio,
            ctime: String(beatmapSet.DateAdded),
          });
        } catch (err) {
          console.error("Error while parsing beatmap: ", err);
        }
      }
    } catch (err) {
      console.error("Error while parsing beatmapset: ", err);
    }
  }

  realm.close();
  console.log(JSON.stringify(Object.fromEntries(songTable)));
  process.exit();
}

parseLazerDatabase().catch(console.error);
