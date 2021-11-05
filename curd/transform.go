package curd

import (
	"go.mongodb.org/mongo-driver/bson"
)

type DBRef struct {
	Ref interface{} `bson:"$ref"`
	ID  interface{} `bson:"$id"`
	DB  interface{} `bson:"$db"`
}

// 对于有引用的字段，$ref 必须放在 $id 之前
func transform(src bson.A) (dest bson.A) {
	for _, value := range src {
		switch v := value.(type) {
		case bson.M:
			i, ok := v["$ref"]
			if ok {
				var d DBRef
				d.Ref = i
				d.ID = v["$id"]
				t, ok := v["$db"]
				if ok {
					d.DB = t
				} else {
					d.DB = ""
				}
				dest = append(dest, d)
			} else {
				dest = src
			}
		default:
			dest = src
		}
	}
	return
}

func format(full bson.M) bson.M {
	for key, value := range full {
		switch v := value.(type) {
		case bson.A:
			full[key] = transform(v)
		default:
			full[key] = value
		}
	}
	return full
}
