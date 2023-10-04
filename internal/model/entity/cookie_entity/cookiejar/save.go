package cookiejar

import "github.com/scriptscat/cloudcat/internal/model/entity/cookie_entity"

func (j *Jar) Import(cookies []*cookie_entity.HttpCookie) {

}

func (j *Jar) Export() map[string][]*cookie_entity.HttpCookie {

	return nil
}
