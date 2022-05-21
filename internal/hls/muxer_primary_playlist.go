package hls

import (
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"github.com/aler9/gortsplib"
)

type muxerPrimaryPlaylist struct {
	fmp4       bool
	videoTrack *gortsplib.TrackH264
	audioTrack *gortsplib.TrackAAC
}

func newMuxerPrimaryPlaylist(
	fmp4 bool,
	videoTrack *gortsplib.TrackH264,
	audioTrack *gortsplib.TrackAAC,
) *muxerPrimaryPlaylist {
	return &muxerPrimaryPlaylist{
		fmp4:       fmp4,
		videoTrack: videoTrack,
		audioTrack: audioTrack,
	}
}

func (p *muxerPrimaryPlaylist) file() *MuxerFileResponse {
	return &MuxerFileResponse{
		Status: http.StatusOK,
		Header: map[string]string{
			"Content-Type": `application/x-mpegURL`,
		},
		Body: &asyncReader{generator: func() []byte {
			var codecs []string

			if p.videoTrack != nil {
				sps := p.videoTrack.SPS()
				if len(sps) >= 4 {
					codecs = append(codecs, "avc1."+hex.EncodeToString(sps[1:4]))
				}
			}

			// https://developer.mozilla.org/en-US/docs/Web/Media/Formats/codecs_parameter
			if p.audioTrack != nil {
				codecs = append(codecs, "mp4a.40."+strconv.FormatInt(int64(p.audioTrack.Type()), 10))
			}

			switch {
			case !p.fmp4:
				return []byte("#EXTM3U\n" +
					"#EXT-X-VERSION:3\n" +
					"\n" +
					"#EXT-X-STREAM-INF:BANDWIDTH=200000,CODECS=\"" + strings.Join(codecs, ",") + "\"\n" +
					"stream.m3u8\n")

			default:
				return []byte("#EXTM3U\n" +
					"#EXT-X-VERSION:7\n" +
					"\n" +
					"#EXT-X-STREAM-INF:BANDWIDTH=200000,CODECS=\"" + strings.Join(codecs, ",") + "\"\n" +
					"stream.m3u8\n" +
					"\n")
			}
		}},
	}
}
