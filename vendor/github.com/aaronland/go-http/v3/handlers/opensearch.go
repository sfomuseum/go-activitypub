package handlers

import (
	"encoding/xml"
)

const DEFAULT_IMAGE_HEIGHT int = 16
const DEFAULT_IMAGE_WIDTH int = 16
const DEFAULT_URL_TYPE string = "text/html"
const DEFAULT_URL_METHOD string = "GET"
const DEFAULT_SEARCHTERMS string = "{searchTerms}"

// https://thenounproject.com/search/?q=globe&i=2115311

const DEFAULT_IMAGE_URI string = "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAAGAAAABgCAYAAADimHc4AAAAAXNSR0IArs4c6QAAAJBlWElmTU0AKgAAAAgABgEGAAMAAAABAAIAAAESAAMAAAABAAEAAAEaAAUAAAABAAAAVgEbAAUAAAABAAAAXgEoAAMAAAABAAIAAIdpAAQAAAABAAAAZgAAAAAAAABIAAAAAQAAAEgAAAABAAOgAQADAAAAAQABAACgAgAEAAAAAQAAAGCgAwAEAAAAAQAAAGAAAAAAEmmPaQAAAAlwSFlzAAALEwAACxMBAJqcGAAAAgtpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IlhNUCBDb3JlIDUuNC4wIj4KICAgPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4KICAgICAgPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIKICAgICAgICAgICAgeG1sbnM6dGlmZj0iaHR0cDovL25zLmFkb2JlLmNvbS90aWZmLzEuMC8iPgogICAgICAgICA8dGlmZjpPcmllbnRhdGlvbj4xPC90aWZmOk9yaWVudGF0aW9uPgogICAgICAgICA8dGlmZjpQaG90b21ldHJpY0ludGVycHJldGF0aW9uPjI8L3RpZmY6UGhvdG9tZXRyaWNJbnRlcnByZXRhdGlvbj4KICAgICAgICAgPHRpZmY6UmVzb2x1dGlvblVuaXQ+MjwvdGlmZjpSZXNvbHV0aW9uVW5pdD4KICAgICAgICAgPHRpZmY6Q29tcHJlc3Npb24+NTwvdGlmZjpDb21wcmVzc2lvbj4KICAgICAgPC9yZGY6RGVzY3JpcHRpb24+CiAgIDwvcmRmOlJERj4KPC94OnhtcG1ldGE+Cs+OiooAAA1TSURBVHgB7Vt7sJVVFb+8xARE3shDr8pDSkXKmBLGIJQssfpD0EjHB4yOZmgjYxOiWanZUI0WmaGCWWg2UqNkkyYhZgok4BiQPFTu5YKggfEQEbjV73f6FvzuYn/nfudy7j3fufdbM7+z115r7dfae6+993ehoiKjzAOZBzIPZB7IPJB5IPNAS/RAqzIedB/0fTDQG+gUoQ3Sx4EqIKMieqAj6vo8MANYAuwE/huDVyHPqAgeOBp1jAfmA/uAOId7eVlNQFsMLG10PDo0FZgEdM7Tue3QrQEYbrgjdgHbgEeBjBrgAcb0+4C9gF/V/4FsOfBD4AtAdyCjInmAu/AbQCiur4X8VqASyKgRPDAUdb4G+BX/EmQXNEJ7qauylNfQa+CNewAetkYrwdwILDBBnvRk6E4FTgR4Fe0AtAcYwt4HNgM8H/4JVAMZRR6gk+YCuup5gPLgzXcp6A/9tQBvRTxstXx9/FbYzwN4sPcCWizxwcTVrQ5bgfzAGI/wYTUB+DPAg1jLNZQ/gHqeAsYBpYwAaL5piTcX3mTUcfcjzx3hqTUEVwLrALXPx9N2OrCpgDI8f/jWaPbEl+xSQB14c8yoPw35Mmer5YzXHcF4z3OANADYAYTsTObTRbA/A2iW1A6jehawQTMEXB4YKe2+D9QCZhtK6fgHAB60pudk0vGGH4juXfCzJW9lfMrX9jcB7r5mRbMwGh3s5MDouHoXOzst05Q8F8txgT6WpWiic+q0wChOh2yjs2tKh4faeh39OSnQ16KLGvMWMAi9ZSxn/Cf9Grgsxx36ORPsc0C3SMTwtCHi+yG1NwLD0luR/BikfSJe7SNRneQU5GyMNeA/jPJ8QxjtBrMlynj70ZCvN8NyS59Hh211cUXZRNg4BoPZJjaM6fzOQ+Lq+wCw8t+iMKIpSE3+pAlj0pfF9mKx+YnIt4PnDY10KbAfsPq5M/sCZUdc6TaIfeD9DaMHZG+IzXvgRwBGvwVj5Xl7amMKpHyIme4akYfY6WL7oBhwF3FlWz0/FR0/gewR3avg/eIR8/Sxx6JL3NI2uLtdFxkSnhU9B/tJsWFYsismU9WxLCfL6o57wMEkR5xUs10bySwZJzoukhNMgXQscACwso+KLvXsNOl4FXh+o1HiVc8GxvQrqgT/mOgfd7ohorO47UzqZI9Cbq+U6VpHW1GxSHT3ON0NomM/r3D6VGa5tXnvNgf7Q5crVh1ylxvF8chrDB7q9BbaaiGf63RxWT0HznNGn0Pe+roTvA81DFum587rBaSabkTvrMOM8W1db58RvY/tNNXd8xdXltkZgNX/nYA+JPq5lLkpYLBK9JOcngtqnegfcfrUZVdLZ692vRstOsZXfzDTfKXY+NBE/dOiH09BAroONjZpswP2/GOQ6RcE9J8VPc8kvltSSWehVzYQbmeuHqUXkDH9L1QR8YNF/z54X55muhqTOmKM1PsiK3HEsMeQxr5xYdiVFOxB+j046/sTB6UpY/Ru7VfaMBnAh+D7B/quKzE0SF5F9Xz4SKCOkOgECM15m0MGkHFizGZiwIZnEVe/TRLrTB1tQI9sEAw3SrOQMd3DqhD+D2IzWeTGctKsjq0mTJDqxNGJvBl5mg6B1f2QV0b5P4nN92JsChb7Q7LgCqICfNqfGPEMPww3RmzjIssgfQUYLnljzzYG6S7A23xM9P8O6EV9GMsbTA+gFcCbj7/C8kVu9Bkwvm3qXgJYlnQJcGuOS8kPV6ytIL5UlTQGm01zSJOeQeqLw/jWh0kaJhghxRYKT3asyzeXrH9TNGhcxQpBp0rrK4QnO1Lyq8DzhuOJ8Z23EdImoCbHHfqSyWwfoF8kfxvpxohPkpwEI4Yg0lvAOzmu7g9tekaiaqShc6Yz5IMiG47rxxFf8mQ7emBhxRzJTjHmMp6bri+FAdLH0tcDeop48Fk9t8XYxIlnStnrY4zuEJvbY2z0NvdmjE1B4mKEIN6bu0St0tlcnUZc2R2jzA6kXN0h0omx1e/t+JHPiAd9IcR+GXUyxqXarvZHzV5HhjcpUiWQ9CpM+yAVYwK6Sc2bhSdrIYP8Gv7EUFeR/0t4ZdtLht+TCqEPxPho4ZXlNywjHZPJmLIehicSd3fcROUMkvwUYwJ0RelKY/u9pBMMU3Gkq5s7JUTtRLhf+CSs2ms9Wlbb1f6oDXkdh50Z3iZxvrEnQLdovrChq5sv5RCp49ShIVsvU3utR+10V2l/1Ia8jiP0ucTb580XYwKOkhb2CU+2reTjHEsT7UetlFFWbSwOqz4fr/Zaj5ZJYkN7HaOOT+tKzMd1JnEFMNwtxnbgmkhXVb5tnWSFqk2hA1f7A9Y5l+rO0LacWe7/o5lMx2eygtJiTIDGfT0P2JGkcXWP9NpPoqnUKeos0+dL1V7r0TLabuitYrZ8Cxjxk8gRUTEmQGNiF9cbvRXxShpHOhBfh5VRx2nYM32+VO01hGgZbVf7ozbk9WanV25vlyhfjAmoQUu2rfuC12seX518PJFOBlSXE0Y/WySjDzkR13lB62pVmzhe7eNWN1/aRqFXMHV0voVS7vzQi5p2iUljY+JCzpDOp6MHApxQpv8ASDwfqoBKgLorAOY96W1iFJQhB+jEnAmb830lefJDRFcJ/gKA93j2yXAOeCOGmVD9Q80A6WrAFpeIS8POjzrDDk1wXXhMdNQ3F9zrxtmgLGe/GLRcKhkpPNkXXL65ZBcVYyDchsWgUahkYVTRKqSnRTwTxs1qwNp6EbxeXZHN0Vj82oJgXf6Kx0PyUznLigoeki9HfJJkFIzsUcgFobculmdcP5sMiO3aWHKC6Kcb0uERzzcNv4GFxhGZNG3Cw5XfSSy89HbNLxHd9U5nWZ4bVp5/xPH0UQhMv9Yr68lzwqwsHenpqxCYfoFXRnn9GvtUjE3BYltxBRd0BbhqdEte7PQPSf5r4EPtLhSbc4U3dqMxSHmltR0l4iDLA9Xu7lz52wJW2t7zAT0X2GSRzxY+Nexl6ImtomWuVx2Q3yH68U7P7IWiZxgLEZ1nbei1MWRrsmFSZrUJJW0DntdJq9fCnJhUcNeavho8y6SOeNfmHds6errroW7hddDp44imjNGMqVaeIccT477pR3llTJ670co8GbAZLfrN4P3u5PmwVWziQihMCiffWOE1HCpB5z1xKFtxs/BkfwTYC3MA+KkUCvEMeVryVwpv7BpjkIYmSNQHWbXT8mZwlTFI5wH6UY6q7wI9yYAYBh/IcSn9OQ394gC44vhAOwVQ4uqx1chzg/ZKfPyYnn8gsZuL2dwk+vtNWE86B3rrE8OkUndk9PLAcKV0DjK1gPUpFDrVPhX8fOnwL12PuOP+Lvr14LuKDfVvip4HttIYZMwZr6giD18jZXxY5Oq2+pa6OnjQbxH9H50+tdnh6JmtOKYjXE8HIa+xfgHybcVGdwmdp7uAtxmrex/4Y6RciOUbxBzMNtuIEVd/3MWA9S4DrOw74PVTCLLppjnonnX+NfDqYPacocD0TB8GzDl0+CbA9N8Gr8T6THeuKgL8pWL7nNPfJzreurj7SO0B3vOtDYYghsayol7ora6uOwK9vxMyGyTT3wEcPOkqwHQ8K4ZQGNG9SE03w4Qx6SNie4vYcFfSsVbPuEjXCSl3pMmZTol0jZLYqit25byObgMujCoeiZRXSMZ3o4Vg+gN28NHJdMwBgAMfBXA3cPew/C7gDIBlGOZIvQG2Q3kIdJ6FqcXgaU+7mUAXgFQN/A2gnAe7hsy7kb8LKFuah57batoKvtKNpBXyswCzSVPKHVr2xFVWBZhjeQ/vERjVdMjscDXbUqX70ZdrA30sWxG39g7AHMpr6LGB0XwJMoYTs4tLOVHPODvm5woWiT5kH1f3JpQbDTQ7Og8j4rXRBr4CPOOxJ37f4X3b7OLS9bDRq+yvXEU6Qe9CV5Ogzt/ARt8krsryz16CIXB7m1PfAD8wZlhfhpx6s02SToU9LxU8OJPYm81K2I8BWgQxzOwFbPAMTRNiRs7bz+XAasDs60v3FGC7DLYXAa2BFkUMR+8B6kxe/zrl8QKvoQ8CWwAtVyi/EeVnAp8AWjQNwOj1RUtH8gDkp+N8xGvrWcAUgHGfBzo/E/iJ4MH7NrAYmA1cB/AykJF4oAP4OYB33l8hGyt2SViGkY5A9yhtcWEliZPibM6HYgPgJ2IpZBMBe8mCzaixPMDVeyewG/ATsRMyhpEvAp2BZkWMqWminujMNOBqgN+BPNVCsBxYAvBVTVQBnCR+K7JbENiMjsQDx6HwDUAhV1DuHF5xbwEyKqIHPo66bgN4JvB240OUz2+HTdlQ2kJQfY7jzhgGcFKGAv0A/u2hN9AVIDEM8YaVURN6gA63ncC/RZQNZfflEk9VNgHZBJTYAyVuPtsB2QSU2AMlbj7bAdkElNgDJW4+2wHZBJTYAyVunn93TSvxdTsJsE8M+frZTpTkb5d8PpZ/4uSnbv6LjYycB36GvH1eaMyUH/pKRtkZUDLX/7/hNH8NLSQENdSNWQhqqOeycpkHMg9kHsg8kHkg80DmgSPywP8AmvVAikZ7npkAAAAASUVORK5CYII="

const NS_MOZ string = "http://www.mozilla.org/2006/browser/search/"
const NS_OPENSEARCH string = "http://a9.com/-/spec/opensearch/1.1/"

type OpenSearchImage struct {
	Height int    `xml:"height,attr"`
	Width  int    `xml:"width,attr"`
	URI    string `xml:",chardata"`
}

type OpenSearchURL struct {
	Type       string                    `xml:"type,attr"`
	Method     string                    `xml:"method,attr"`
	Template   string                    `xml:"template,attr"`
	Parameters []*OpenSearchURLParameter `xml:"Param"`
}

type OpenSearchURLParameter struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type OpenSearchDescription struct {
	XMLName       xml.Name         `xml:"OpenSearchDescription"`
	NSMoz         string           `xml:"xmlns:moz,attr"`
	InputEncoding string           `xml:"InputEncoding"`
	NSOpenSearch  string           `xml:"xmlns,attr"`
	ShortName     string           `xml:"ShortName"`
	Description   string           `xml:"Description"`
	Image         *OpenSearchImage `xml:"Image"`
	URL           *OpenSearchURL   `xml:"Url"`
	SearchForm    string           `xml:"moz:searchForm"`
}

func (d *OpenSearchDescription) Marshal() ([]byte, error) {

	enc, err := xml.Marshal(d)

	if err != nil {
		return nil, err
	}

	body := []byte(xml.Header)
	body = append(body, enc...)

	return body, nil
}
