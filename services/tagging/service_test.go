package tagging

// func TestCreate(t *testing.T) {
// 	t.Run("Valid", func(t *testing.T) {
// 		store := &tag.RepositoryMock{
// 			SaveFunc: func(tags domain.Tags) error {
// 				return nil
// 			},
// 			GetByNameFunc: func(names string) (domain.Tags, error) {
// 				return domain.Tags{}, nil
// 			},
// 		}

// 		svc := New(store)
// 		_, err := svc.Create("tag 1")

// 		is := is.New(t)
// 		is.Equal(err, nil)
// 		is.Equal(len(store.SaveCalls()), 1)
// 	})

// 	t.Run("Invalid: blank tag name", func(t *testing.T) {
// 		store := &tag.RepositoryMock{
// 			SaveFunc: func(tags domain.Tags) error {
// 				return nil
// 			},
// 		}

// 		svc := New(store)
// 		_, err := svc.Create("")

// 		is := is.New(t)
// 		is.Equal(err, domain.ErrBlankTag)
// 		is.Equal(len(store.SaveCalls()), 0)
// 	})

// 	t.Run("Invalid: tag name too long", func(t *testing.T) {
// 		store := &tag.RepositoryMock{
// 			SaveFunc: func(tags domain.Tags) error {
// 				return nil
// 			},
// 		}

// 		svc := New(store)
// 		_, err := svc.Create("Lorem Ipsum is simply dummy text of the printing and typesetting industry.")
// 		is := is.New(t)
// 		is.True(err != nil)
// 		is.Equal(len(store.SaveCalls()), 0)
// 	})
// }

// func TestUpdate(t *testing.T) {
// 	t.Run("Valid", func(t *testing.T) {
// 		tg, _ := domain.NewTags("tag 1")

// 		store := &tag.RepositoryMock{
// 			GetByIdFunc: func(id string) (*domain.Tags, error) {
// 				return tg, nil
// 			},
// 			UpdateFunc: func(tagsIn domain.Tags) error {
// 				tg = &tagsIn
// 				return nil
// 			},
// 		}

// 		svc := New(store)
// 		err := svc.Update(tg.ID.String(), "tag 2")
// 		is := is.New(t)
// 		is.Equal(err, nil)
// 		is.Equal(len(store.UpdateCalls()), 1)
// 	})

// 	t.Run("Invalid: Not found", func(t *testing.T) {
// 		store := &tag.RepositoryMock{
// 			GetByIdFunc: func(id string) (*domain.Tags, error) {
// 				return nil, sql.ErrNoRows
// 			},
// 			UpdateFunc: func(tagsIn domain.Tags) error {
// 				return nil
// 			},
// 		}

// 		svc := New(store)
// 		err := svc.Update("tag item is not found", "tag 2")
// 		is := is.New(t)
// 		is.Equal(err, sql.ErrNoRows)
// 		is.Equal(len(store.UpdateCalls()), 0)
// 	})
// }

// func TestDelete(t *testing.T) {
// 	t.Run("Valid", func(t *testing.T) {
// 		store := &tag.RepositoryMock{
// 			GetByIdFunc: func(s string) (*domain.Tags, error) {
// 				return nil, nil
// 			},
// 			DeleteFunc: func(s string) error {
// 				return nil
// 			},
// 		}

// 		svc := New(store)
// 		is := is.New(t)

// 		err := svc.Delete(uuid.New().String())
// 		is.Equal(err, nil)
// 		is.Equal(len(store.DeleteCalls()), 1)
// 	})
// }
