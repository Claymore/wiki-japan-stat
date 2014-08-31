#!/usr/bin/python
# -*- coding: utf-8 -*-
import csv
import codecs
import ConfigParser

def unicode_csv_reader(utf8_data, dialect=csv.excel, **kwargs):
    reader = csv.reader(utf8_data, dialect=dialect, **kwargs)
    for row in reader:
        yield [unicode(cell, 'utf-8') for cell in row]

config = ConfigParser.RawConfigParser()
config.readfp(codecs.open("prefectures.cfg", "r", "utf8"))

for code in config.sections():
    ja_prefecture_name = u''
    filename = 'data/%s.csv' % (code)
    template_name = 'templates/%s.txt' % (code)

    name_column_id = config.getint(code, 'name_column_id')
    population_column_id = config.getint(code, 'population_column_id')
    population_column_id = config.getint(code, 'population_column_id')
    name_column_2_id = 0
    if config.has_option(code, 'name_column_2_id'):
        name_column_2_id = config.getint(code, 'name_column_2_id')
    population_column_2_id = 0
    if config.has_option(code, 'population_column_2_id'):
        population_column_2_id = config.getint(code, 'population_column_2_id')
    publication_date = config.get(code, 'publication_date')
    has_stop_after = config.has_option(code, 'stop_after')
    stop_after = 0
    if has_stop_after:
        stop_after = config.getint(code, 'stop_after')
    url = config.get(code, 'url')
    title = config.get(code, 'title')
    publisher = config.get(code, 'publisher')
    access_date = config.get(code, 'access_date')
    description = config.get(code, 'description')
    data_date = config.get(code, 'data_date')
    total_name = config.get(code, 'total_name')

    codes = dict()
    with open('codes.csv', 'rb') as f:
        reader = unicode_csv_reader(f, delimiter='\t')
        for row in reader:
            if row[0].startswith(code[3:]):
                codes[row[2]] = row[0]
                ja_prefecture_name = row[1]

    with open(filename, 'rb') as fr:
        fw = codecs.open(template_name, 'w', 'utf-8')
        fw.write(u'{{#switch:{{{1}}}\n  |источник = <ref>{{cite web|date=%s|url=%s|title=%s|publisher=%s|accessdate=%s|lang=jp|description=%s}}</ref>\n' % (publication_date, url, title, publisher, access_date, description))
        fw.write(u'  |дата    = %s\n' % (data_date))
        municipalities = 0
        reader = unicode_csv_reader(fr, delimiter=';')
        line_number = 1
        for row in reader:
            if has_stop_after and line_number > stop_after:
                break
            line_number = line_number + 1
            name = row[name_column_id].strip().replace(' ', '')
            population = u''.join(row[population_column_id].split())
            if population.isnumeric():
                if name == total_name:
                    fw.write(u'  |%s   = %s\n' % (code, population))
                if name.endswith(u'市') or name.endswith(u'町') or name.endswith(u'村'):
                    if name in codes:
                        fw.write('  |%s = %7s <!-- %s -->\n' % (codes[name], population, name))
                        municipalities = municipalities + 1
                    else:
                        print '[warning] Skipping %s %s' % (code, name)
            if name_column_2_id != 0:
                name = row[name_column_2_id].strip().replace(' ', '')
                population = u''.join(row[population_column_2_id].split())
                if population.isnumeric():
                    if name == total_name:
                        fw.write(u'  |%s   = %s\n' % (code, population))
                    if name.endswith(u'市') or name.endswith(u'町') or name.endswith(u'村'):
                        if name in codes:
                            fw.write('  |%s = %7s <!-- %s -->\n' % (codes[name], population, name))
                            municipalities = municipalities + 1
                        else:
                            print '[warning] Skipping %s %s' % (code, name)
        fw.write(u'  |муниципалитетов = %d\n' % (municipalities))
        fw.write(u'  |#default = <span style="color:red;">Неверный параметр</span><includeonly>[[Категория:Википедия:Населённые пункты Японии с ошибочным параметром шаблона]]</includeonly>\n')
        fw.write(u'}}<noinclude>\n')
        fw.write(u'{{doc|Население административной единицы Японии/doc}}\n\n')
        fw.write(u'[[Категория:Шаблоны:Население административных единиц Японии|%s]]\n\n</noinclude>' % (code))
        fw.close()
