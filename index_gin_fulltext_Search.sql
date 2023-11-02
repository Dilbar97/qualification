create table ts(doc text, doc_tsv tsvector);

insert into ts(doc) values
                        ('Во поле береза стояла'),  ('Во поле кудрявая стояла'),
                        ('Люли, люли, стояла'),     ('Люли, люли, стояла'),
                        ('Некому березу заломати'), ('Некому кудряву заломати'),
                        ('Люли, люли, заломати'),   ('Люли, люли, заломати'),
                        ('Я пойду погуляю'),        ('Белую березу заломаю'),
                        ('Люли, люли, заломаю'),    ('Люли, люли, заломаю');

set default_text_search_config = russian;

update ts set doc_tsv = to_tsvector(doc) WHERE doc_tsv IS null;

create index on ts using gin(doc_tsv);

select ctid, doc, doc_tsv from ts;

select doc from ts where doc_tsv @@ to_tsquery('стояла & кудрявая');