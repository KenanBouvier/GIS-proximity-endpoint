create function get_domain_name(url varchar)
returns varchar
language plpgsql
as
$$
declare
   web varchar;
begin
   select
   into web
   		substring(url from '(?:.*://)?(?:www\.)?([^/?]*)');
	return web;
end;
$$;

-- Run pl/pgsql function with ` SELECT get_domain_name('http://test.com/index.php'); ` -> 'test.com'
