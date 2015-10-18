package codegen

// GetNewGodepsJsonTemplate => builds the template string for generating the Godep json file
func GetNewGodepsJsonTemplate(cc CodeConfig) string {
	return `
{
	"ImportPath": "github.com/titanium-io/titanium-api",
	"GoVersion": "go1.4.2",
	"Deps": [
		{
			"ImportPath": "github.com/BurntSushi/toml",
			"Comment": "v0.1.0-21-g056c9bc",
			"Rev": "056c9bc7be7190eaa7715723883caffa5f8fa3e4"
		},
		{
			"ImportPath": "github.com/Sirupsen/logrus",
			"Comment": "v0.8.6-1-g8bca266",
			"Rev": "8bca2664072173a3c71db4c28ca8d304079b1787"
		},
		{
			"ImportPath": "github.com/cenkalti/backoff",
			"Rev": "6c45d6bc1e78d94431dff8fc28a99f20bafa355a"
		},
		{
			"ImportPath": "github.com/codegangsta/inject",
			"Comment": "v1.0-rc1-10-g33e0aa1",
			"Rev": "33e0aa1cb7c019ccc3fbe049a8262a6403d30504"
		},
		{
			"ImportPath": "github.com/dancannon/gorethink",
			"Comment": "v1.x.x",
			"Rev": "8aca6ba2cc6e873299617d730fac0d7f6593113a"
		},
		{
			"ImportPath": "github.com/gedex/inflector",
			"Rev": "8c0e57904488c554ab26caec525db5c92b23f051"
		},
		{
			"ImportPath": "github.com/go-martini/martini",
			"Comment": "v1.0-166-g9987fc5",
			"Rev": "9987fc59f060a063247d22472c2fb48654bac5f6"
		},
		{
			"ImportPath": "github.com/golang/protobuf/proto",
			"Rev": "68c687dc49948540b356a6b47931c9be4fcd0245"
		},
		{
			"ImportPath": "github.com/kr/pretty",
			"Comment": "go.weekly.2011-12-22-27-ge6ac2fc",
			"Rev": "e6ac2fc51e89a3249e82157fa0bb7a18ef9dd5bb"
		},
		{
			"ImportPath": "github.com/kr/text",
			"Rev": "bb797dc4fb8320488f47bf11de07a733d7233e1f"
		},
		{
			"ImportPath": "github.com/magiconair/properties",
			"Comment": "v1.5.5",
			"Rev": "337395e44efb4affc8bf11971e22acbd4e27a32d"
		},
		{
			"ImportPath": "github.com/manyminds/api2go/jsonapi",
			"Comment": "0.5-2-g72259cb",
			"Rev": "72259cbf8cd22938ae3428c0e8db2cf49f3f7049"
		},
		{
			"ImportPath": "github.com/martini-contrib/cors",
			"Rev": "553b9208d353a39b0850c02355f478ba020c86d7"
		},
		{
			"ImportPath": "github.com/martini-contrib/gorelic",
			"Rev": "f8b0843aa3ab66d8734aa7d9acf399bdc197c65d"
		},
		{
			"ImportPath": "github.com/martini-contrib/render",
			"Rev": "ec18f8345a1181146728238980606fb1d6f40e8c"
		},
		{
			"ImportPath": "github.com/mitchellh/mapstructure",
			"Rev": "281073eb9eb092240d33ef253c404f1cca550309"
		},
		{
			"ImportPath": "github.com/obieq/gas",
			"Rev": "c3d0cdef2ab9d6050a5327d88c44e04aa598cb3a"
		},
		{
			"ImportPath": "github.com/obieq/goar",
			"Rev": "6572645b80930c158cb5ff3ac3a5000838d85002"
		},
		{
			"ImportPath": "github.com/obieq/goar-validations",
			"Rev": "aba3dc402a8fd7f92a431be938b723b9c08b82d6"
		},
		{
			"ImportPath": "github.com/obieq/gson-api",
			"Rev": "a9acf85d1b32f033aad2798b40bdd73259447555"
		},
		{
			"ImportPath": "github.com/oxtoacart/bpool",
			"Rev": "4e1c5567d7c2dd59fa4c7c83d34c2f3528b025d6"
		},
		{
			"ImportPath": "github.com/spf13/cast",
			"Rev": "ee815aaf958c707ad07547cd62150d973710f747"
		},
		{
			"ImportPath": "github.com/spf13/jwalterweatherman",
			"Rev": "3d60171a64319ef63c78bd45bd60e6eab1e75f8b"
		},
		{
			"ImportPath": "github.com/spf13/pflag",
			"Rev": "b084184666e02084b8ccb9b704bf0d79c466eb1d"
		},
		{
			"ImportPath": "github.com/spf13/viper",
			"Rev": "1967d93db724f4a5c0e101307e96d82ff520a067"
		},
		{
			"ImportPath": "github.com/twinj/uuid",
			"Rev": "70cac2bcd273ef6a371bb96cde363d28b68734c3"
		},
		{
			"ImportPath": "github.com/yvasiyarov/go-metrics",
			"Rev": "c25f46c4b94079672242ec48a545e7ca9ebe3aec"
		},
		{
			"ImportPath": "github.com/yvasiyarov/gorelic",
			"Comment": "v0.0.6-30-g9ad1745",
			"Rev": "9ad1745d72a33d259bbc1c168b28cd728ad9e557"
		},
		{
			"ImportPath": "github.com/yvasiyarov/newrelic_platform_go",
			"Rev": "db03d6dfd965ebc22dbfd910567c641b326ff055"
		},
		{
			"ImportPath": "gopkg.in/guregu/null.v3",
			"Comment": "v3",
			"Rev": "a9db3ac26fcd2d70230a8dc6286d0c3694700003"
		},
		{
			"ImportPath": "gopkg.in/yaml.v2",
			"Rev": "53feefa2559fb8dfa8d81baad31be332c97d6c77"
		}
	]
}`
}
