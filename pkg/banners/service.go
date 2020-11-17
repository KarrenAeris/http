package banners

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"sync"
)

//Service представляет собой сервис по управлению баннерами.
type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

// NewService создаёт сервис.
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

//Banner представляет собой баннер
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

var sID int64 = 0

// All возвращает все существующие баннеры.
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	//возвращаем все баннеры, если их нет, возвращаем []
	return s.items, nil
}

// ByID возвращает баннеры по идентификатору.
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.items {
		//если ID элемента равен ID из параметра, то мы нашли баннер
		if v.ID == id {
			//вернем баннер и ошибку nil
			return v, nil
		}
	}

	return nil, errors.New("item not found")
}

//Save сохраяет/обновляет баннер.
func (s *Service) Save(ctx context.Context, item *Banner, file multipart.File) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if item.ID == 0 {
		sID++
		item.ID = sID

		//если файл пришел, то сохраняем его под нужным именем, например: сейчас там только расширение (jpg), а мы его изменим на (2.jpg)
		if item.Image != "" {
			//генерируем имя файла, например: если ID равно 2 и раширение файла jpg , то получим 2.jpg
			item.Image = fmt.Sprint(item.ID) + "." + item.Image
			//вызываем фукции для загрузки файла на сервер и передаем ему файл и path, где нужно сохранить файл  ./web/banners/2.jpg
			err := uploadFile(file, "./web/banners/"+item.Image)
			//если при сохранении произошла какая-нибудь ошибка, то возвращаем ошибку
			if err != nil {
				return nil, err
			}
		}

		s.items = append(s.items, item)
		return item, nil
	}
	for k, v := range s.items {
		if v.ID == item.ID {

			//если файл пришел, то сохраняем его под нужным именем, например: сейчас там только расширение (jpg), а мы его изменим на (2.jpg)
			if item.Image != "" {
				//генерируем имя файла, например: если ID равно 2 и раширение файла jpg , то получим 2.jpg
				item.Image = fmt.Sprint(item.ID) + "." + item.Image
				//вызываем фукции для загрузки файла на сервер и передаем ему файл и path, где нужно сохранить файл  ./web/banners/2.jpg
				err := uploadFile(file, "./web/banners/"+item.Image)
				//если при сохранении произошла какая-нибудь ошибка, то возвращаем ошибку
				if err != nil {
					return nil, err
				}
			} else {
				//если файл не пришел, то просто поставим полученное значение в поле Image
				item.Image = s.items[k].Image
			}

			s.items[k] = item
			return item, nil
		}
	}
	return nil, errors.New("item not found")
}

//RemoveByID удаляет баннер по идентификатору.
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for k, v := range s.items {
		if v.ID == id {
			s.items = append(s.items[:k], s.items[k+1:]...)
			return v, nil
		}
	}

	return nil, errors.New("item not found")
}

//эта функция сохраняет файл в сервере в заданной папке path и возвращает nil, если все успешно, -  или error, если есть ошибка
func uploadFile(file multipart.File, path string) error {
	//прочитаем весь файл и получаем слайс из байтов
	var data, err = ioutil.ReadAll(file)
	//если не удалось прочитать, вернем ошибку
	if err != nil {
		return errors.New("not readble data")
	}

	//записываем файл в заданной папке с публичными правами
	err = ioutil.WriteFile(path, data, 0666)

	//если не удалось записать файл, вернем ошибку
	if err != nil {
		return errors.New("not saved from folder ")
	}

	return nil
}
