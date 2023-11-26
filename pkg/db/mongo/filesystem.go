package mongo

import (
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/zunkk/go-project-startup/pkg/reqctx"
)

type File struct {
	ID         string
	Name       string
	Length     int64
	UploadDate time.Time
	Metadata   string
	Data       io.Reader
}

type Metadata struct {
	Content string `json:"content" bson:"content"`
}

type FileSystemDao struct {
	db   *DB
	fsDB *mongo.Database
}

func NewFileSystem(db *DB) *FileSystemDao {
	d := &FileSystemDao{
		db:   db,
		fsDB: db.Client.Database("gridfs"),
	}
	return d
}

func (d *FileSystemDao) getBucket(bucketName string) (*gridfs.Bucket, error) {
	bucketOptions := options.GridFSBucket().SetName(bucketName)
	return gridfs.NewBucket(d.fsDB, bucketOptions)
}

func (d *FileSystemDao) Upload(ctx *reqctx.ReqCtx, bucketName string, fileName string, data io.Reader, metadata string) (id string, err error) {
	bucket, err := d.getBucket(bucketName)
	if err != nil {
		return "", err
	}
	id = fmt.Sprintf("%d", ctx.RequestID)

	opts := options.GridFSUpload().SetMetadata(&Metadata{
		Content: metadata,
	})
	if err := bucket.UploadFromStreamWithID(id, fileName, data, opts); err != nil {
		return "", err
	}
	return id, nil
}

func (d *FileSystemDao) Download(ctx *reqctx.ReqCtx, bucketName string, id string) (*File, error) {
	bucket, err := d.getBucket(bucketName)
	if err != nil {
		return nil, err
	}

	ds, err := bucket.OpenDownloadStream(id)
	if err != nil {
		return nil, err
	}
	f := ds.GetFile()
	id, _ = f.ID.(string)
	var metadata Metadata
	if err := bson.Unmarshal(f.Metadata, &metadata); err != nil {
		_ = ds.Close()
		return nil, err
	}

	return &File{
		ID:         id,
		Name:       f.Name,
		Length:     f.Length,
		UploadDate: f.UploadDate,
		Metadata:   metadata.Content,
		Data:       ds,
	}, nil
}
