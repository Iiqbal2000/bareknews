INSERT INTO tags (id, name, slug) values
  ("9df76a7f-3d44-4fcf-9222-df65f0774bfb", "Go", "go"),
  ("57cb6822-a459-4d4c-9709-7d0820dc441b", "NodeJS", "nodejs");

INSERT INTO news (id, title, slug, status, body, date_created, date_updated) values 
  ("fdc76bfc-5aa1-4096-a9d3-39719610f987", "Post 2021", "post-2021", "draft", "Kjsjs jji susw", 1257894000, 1257894000),
  ("69da4e7a-f548-4724-a77f-c50e10f08ebc", "Post 2022", "post-2022", "publish", "Kjsjs jji susw1q", 1257894001, 1257894001);

INSERT INTO news_tags(newsID, tagsID) values
  ("fdc76bfc-5aa1-4096-a9d3-39719610f987", "57cb6822-a459-4d4c-9709-7d0820dc441b"),
  ("fdc76bfc-5aa1-4096-a9d3-39719610f987", "9df76a7f-3d44-4fcf-9222-df65f0774bfb")