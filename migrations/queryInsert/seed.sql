INSERT INTO categories (name) VALUES
('Favorite Product'),
('Coffee'),
('Non Coffe'),
('Foods'),
('Add-on');

INSERT INTO variants (name, additional_price) VALUES
('Ice', 7000),
('Hot', 0),
('Food', 0);

INSERT INTO sizes(name, additional_price) VALUES
('R', 0),
('L', 5000),
('XL', 10000),
('250 gr', 50000),
('500 gr', 75000);


INSERT INTO products (title, description, base_price, category_id, stock, is_favorite) VALUES
('Espresso Shot', 'Strong and concentrated espresso, perfect to energize your day', 15000, 1, 40, TRUE),
('Caramel Macchiato', 'Espresso with steamed milk, topped with caramel drizzle', 35000, 1, 25, FALSE),
('Mocha Coffee', 'Rich espresso blended with chocolate and milk', 30000, 1, 20, TRUE),
('Iced Coffee Blend', 'Iced coffee blended with ice and milk for a smooth taste', 35000, 1, 25, FALSE),
('Vanilla Latte', 'Espresso with steamed milk and vanilla syrup', 30000, 1, 20, TRUE),
('Iced Americano', 'Cold espresso diluted with water for a refreshing taste', 25000, 1, 30, FALSE),
('Affogato', 'Vanilla ice cream topped with hot espresso', 35000, 1, 15, TRUE),
('Avocado Coffee', 'Unique coffee with creamy avocado blended in', 35000, 1, 15, FALSE),

('Iced Coffee with Milk', 'Smooth iced coffee blended with fresh milk, lightly sweetened', 25000, 2, 50, TRUE),
('Traditional Black Coffee', 'Authentic Indonesian coffee brewed strong and bold, served hot', 20000, 2, 40, FALSE),
('Cappuccino', 'Espresso topped with rich frothed milk and a hint of cocoa', 30000, 2, 30, TRUE),
('Latte', 'Creamy espresso latte with silky steamed milk', 30000, 2, 25, FALSE),
('Matcha Latte', 'Smooth green tea latte made from premium matcha powder', 30000, 2, 20, TRUE),
('Milo Dinosaur', 'Iced Milo drink with extra Milo powder on top', 25000, 2, 20, FALSE),
('Civet Coffee', 'Premium Indonesian Luwak coffee, rich and aromatic', 100000, 2, 10, TRUE),
('Iced Coffee with Jelly', 'Iced coffee served with sweet coffee jelly', 30000, 2, 20, FALSE),

('Iced Sweet Tea', 'Refreshing iced tea lightly sweetened for a cool thirst-quencher', 15000, 3, 40, TRUE),
('Hot Tea', 'Premium hot tea brewed to perfection', 12000, 3, 50, FALSE),
('Chocolate Milk', 'Rich chocolate blended with fresh cold milk', 20000, 3, 30, TRUE),
('Vanilla Milkshake', 'Creamy vanilla milkshake topped with whipped cream', 25000, 3, 25, FALSE),
('Banana Smoothie', 'Fresh bananas blended with milk for a naturally sweet drink', 25000, 3, 30, TRUE),
('Fresh Mango Juice', 'Pure mango juice made from ripe, sweet mangoes', 20000, 3, 25, FALSE),
('Honey Lemon Drink', 'Refreshing lemon drink sweetened with natural honey', 20000, 3, 25, TRUE),
('Iced Orange Juice', 'Freshly squeezed orange juice served over ice', 18000, 3, 35, FALSE),
('Oreo Milkshake', 'Chocolate milkshake blended with Oreo cookies', 30000, 3, 20, TRUE),
('Iced Chocolate', 'Cold chocolate drink with milk and chocolate syrup', 25000, 3, 30, FALSE),
('Pulled Tea', 'Traditional Malaysian pulled tea with creamy froth', 22000, 3, 30, TRUE),

('Grilled Chicken Sandwich', 'Toasted sandwich with grilled chicken, fresh lettuce, and mayo', 35000, 4, 20, FALSE),
('Beef Burger', 'Juicy beef patty with fresh vegetables, cheese, and sauce', 45000, 4, 20, TRUE),
('French Fries', 'Crispy golden fries, lightly salted', 20000, 4, 50, FALSE),
('Cinnamon Roll', 'Soft pastry rolled with cinnamon and sugar, drizzled with icing', 15000, 4, 30, TRUE),
('Cheesecake', 'Creamy baked cheesecake with a buttery crust', 25000, 4, 15, FALSE),
('Butter Croissant', 'Flaky and buttery croissant, perfect for breakfast', 20000, 4, 30, TRUE),
('Egg Tart', 'Sweet egg custard in a flaky pastry shell', 15000, 4, 25, FALSE),
('Glazed Donut', 'Soft donut coated with sweet sugar glaze', 12000, 4, 40, TRUE),
('Spicy Chicken Wings', 'Crispy chicken wings tossed in spicy sauce', 40000, 4, 20, FALSE),
('Tuna Sandwich', 'Fresh tuna salad with vegetables in toasted bread', 35000, 4, 25, TRUE),
('Soup of the Day', 'Chef’s special soup, freshly made daily', 30000, 4, 15, FALSE),
('Chocolate Brownies', 'Rich chocolate brownies, soft and fudgy', 20000, 4, 25, TRUE),
('Mixed Fruit Salad', 'Fresh seasonal fruits served in a light syrup', 25000, 4, 20, FALSE),

('Extra Espresso Shot', 'Add an extra shot of strong espresso to your drink', 10000, 5, 50, TRUE),
('Whipped Cream Topping', 'Light and fluffy whipped cream topping', 5000, 5, 50, FALSE),
('Caramel Syrup', 'Sweet caramel syrup to enhance your drink', 5000, 5, 50, TRUE),
('Vanilla Syrup', 'Smooth vanilla syrup to add flavor', 5000, 5, 50, FALSE),
('Chocolate Syrup', 'Rich chocolate syrup to top your drink', 5000, 5, 50, TRUE),
('Extra Ice', 'Add extra ice cubes to your drink', 2000, 5, 100, FALSE),
('Brown Sugar Boba', 'Chewy brown sugar boba pearls for your drink', 20000, 5, 30, TRUE),
('Cheese Topping', 'Add creamy cheese topping', 20000, 5, 30, FALSE),
('Hazelnut Syrup', 'Add hazelnut flavor syrup to your drink', 20000, 5, 50, TRUE),
('Chocolate Chips', 'Add chocolate chips for extra sweetness', 21000, 5, 50, FALSE);


UPDATE products SET rating = 4.3 WHERE id = 1;
UPDATE products SET rating = 4.7 WHERE id = 2;
UPDATE products SET rating = 4.5 WHERE id = 3;
UPDATE products SET rating = 4.2 WHERE id = 4;
UPDATE products SET rating = 4.8 WHERE id = 5;
UPDATE products SET rating = 4.4 WHERE id = 6;
UPDATE products SET rating = 4.6 WHERE id = 7;
UPDATE products SET rating = 4.9 WHERE id = 8;
UPDATE products SET rating = 4.1 WHERE id = 9;
UPDATE products SET rating = 4.5 WHERE id = 10;

UPDATE products SET rating = 4.2 WHERE id = 11;
UPDATE products SET rating = 4.7 WHERE id = 12;
UPDATE products SET rating = 4.3 WHERE id = 13;
UPDATE products SET rating = 4.6 WHERE id = 14;
UPDATE products SET rating = 4.9 WHERE id = 15;
UPDATE products SET rating = 4.4 WHERE id = 16;
UPDATE products SET rating = 4.5 WHERE id = 17;
UPDATE products SET rating = 4.8 WHERE id = 18;
UPDATE products SET rating = 4.3 WHERE id = 19;
UPDATE products SET rating = 4.7 WHERE id = 20;

UPDATE products SET rating = 4.6 WHERE id = 21;
UPDATE products SET rating = 4.9 WHERE id = 22;
UPDATE products SET rating = 4.2 WHERE id = 23;
UPDATE products SET rating = 4.5 WHERE id = 24;
UPDATE products SET rating = 4.3 WHERE id = 25;
UPDATE products SET rating = 4.8 WHERE id = 26;
UPDATE products SET rating = 4.4 WHERE id = 27;
UPDATE products SET rating = 4.7 WHERE id = 28;
UPDATE products SET rating = 4.5 WHERE id = 29;
UPDATE products SET rating = 4.9 WHERE id = 30;

UPDATE products SET rating = 4.1 WHERE id = 31;
UPDATE products SET rating = 4.4 WHERE id = 32;
UPDATE products SET rating = 4.6 WHERE id = 33;
UPDATE products SET rating = 4.3 WHERE id = 34;
UPDATE products SET rating = 4.8 WHERE id = 35;
UPDATE products SET rating = 4.2 WHERE id = 36;
UPDATE products SET rating = 4.5 WHERE id = 37;
UPDATE products SET rating = 4.7 WHERE id = 38;
UPDATE products SET rating = 4.9 WHERE id = 39;
UPDATE products SET rating = 4.3 WHERE id = 40;

UPDATE products SET rating = 4.6 WHERE id = 41;
UPDATE products SET rating = 4.2 WHERE id = 42;
UPDATE products SET rating = 4.8 WHERE id = 43;
UPDATE products SET rating = 4.5 WHERE id = 44;
UPDATE products SET rating = 4.7 WHERE id = 45;
UPDATE products SET rating = 4.3 WHERE id = 46;
UPDATE products SET rating = 4.9 WHERE id = 47;
UPDATE products SET rating = 4.4 WHERE id = 48;
UPDATE products SET rating = 4.6 WHERE id = 49;
UPDATE products SET rating = 4.8 WHERE id = 50;


INSERT INTO product_images (product_id, image) VALUES
(1, 'https://images.unsplash.com/photo-1646257861487-60fa89bef25f?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTJ8fEVzcHJlc3NvJTIwU2hvdHxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(1, 'https://images.unsplash.com/photo-1601390483714-955fd3066695?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NHx8RXNwcmVzc28lMjBTaG90fGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(1, 'https://images.unsplash.com/photo-1705952285570-113e76f63fb0?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8RXNwcmVzc28lMjBTaG90fGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(1, 'https://cdn.pixabay.com/photo/2020/09/13/04/13/coffee-5567269_1280.jpg'),
(2, 'https://cdn.pixabay.com/photo/2016/11/02/05/17/coffee-1790576_640.jpg'),
(2, 'https://images.unsplash.com/photo-1570517130750-10c67ffdde09?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8Q2FyYW1lbCUyME1hY2NoaWF0b3xlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(2, 'https://images.unsplash.com/photo-1517701550927-30cf4ba1dba5?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Nnx8Q2FyYW1lbCUyME1hY2NoaWF0b3xlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(2, 'https://images.unsplash.com/photo-1579888071069-c107a6f79d82?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NHx8Q2FyYW1lbCUyME1hY2NoaWF0b3xlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(3, 'https://cdn.pixabay.com/photo/2020/03/07/05/18/coffee-4908764_640.jpg'),
(3, 'https://cdn.pixabay.com/photo/2023/08/07/03/59/coffee-8174279_640.jpg'),
(3, 'https://cdn.pixabay.com/photo/2016/08/09/13/21/coffee-1580595_640.jpg'),
(3, 'https://cdn.pixabay.com/photo/2017/03/17/10/29/coffee-2151200_640.jpg'),
(4, 'https://cdn.pixabay.com/photo/2021/12/17/09/13/iced-coffee-6876063_640.jpg'),
(4, 'https://cdn.pixabay.com/photo/2018/10/19/16/45/coffee-3759003_640.jpg'),
(4, 'https://cdn.pixabay.com/photo/2020/09/07/13/41/cappuccino-5551909_640.jpg'),
(4, 'https://cdn.pixabay.com/photo/2023/01/28/11/58/iced-coffee-7750749_640.jpg'),
(5, 'https://images.unsplash.com/photo-1504194472231-5a5294eddc43?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8VmFuaWxsYSUyMExhdHRlfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(5, 'https://images.unsplash.com/photo-1504194472231-5a5294eddc43?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8VmFuaWxsYSUyMExhdHRlfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(6, 'https://cdn.pixabay.com/photo/2020/05/28/21/25/coffee-5232766_640.jpg'),
(6, 'https://cdn.pixabay.com/photo/2022/04/05/07/46/iced-coffee-7113044_640.jpg'),
(6, 'https://cdn.pixabay.com/photo/2021/08/24/04/52/coffee-6569483_640.jpg'),
(6, 'https://cdn.pixabay.com/photo/2019/04/14/17/09/americano-4127162_640.jpg'),
(7, 'https://cdn.pixabay.com/photo/2015/02/04/09/06/affogato-623516_640.jpg'),
(7, 'https://cdn.pixabay.com/photo/2017/08/05/17/02/espresso-2584668_640.jpg'),
(8, 'https://www.nescafe.com/id/sites/default/files/2025-08/Hero%20banner%20desktop.jpg'),
(9, 'https://images.unsplash.com/photo-1721466810960-e6a5b252d9a4?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8SWNlZCUyMENvZmZlZSUyMHdpdGglMjBNaWxrfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(9, 'https://images.unsplash.com/photo-1629688825560-917b9727c15b?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8SWNlZCUyMENvZmZlZSUyMHdpdGglMjBNaWxrfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(9, 'https://images.unsplash.com/photo-1709689242523-c6d8027c7499?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NHx8SWNlZCUyMENvZmZlZSUyMHdpdGglMjBNaWxrfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(9, 'https://images.unsplash.com/photo-1586551199331-91872af420c9?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Nnx8SWNlZCUyMENvZmZlZSUyMHdpdGglMjBNaWxrfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(10, 'https://cdn.pixabay.com/photo/2017/10/27/11/19/coffee-2893970_640.jpg'),
(10, 'https://cdn.pixabay.com/photo/2015/01/03/22/03/cup-of-coffee-587873_640.jpg'),
(10, 'https://images.unsplash.com/photo-1642257980462-f720b2c01346?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Nnx8dHJhZGl0aW9uYWwlMjBibGFjayUyMGNvZmZlZXxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(10, 'https://images.unsplash.com/photo-1677327109315-04debf6c5f21?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8N3x8dHJhZGl0aW9uYWwlMjBibGFjayUyMGNvZmZlZXxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(11, 'https://images.unsplash.com/photo-1534687941688-651ccaafbff8?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTF8fENhcHB1Y2Npbm98ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(11, 'https://images.unsplash.com/photo-1503240778100-fd245e17a273?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTJ8fENhcHB1Y2Npbm98ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(11, 'https://images.unsplash.com/photo-1529892485617-25f63cd7b1e9?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTh8fENhcHB1Y2Npbm98ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(11, 'https://images.unsplash.com/photo-1486122151631-4b5aaa3ac70d?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTl8fENhcHB1Y2Npbm98ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(12, 'https://cdn.pixabay.com/photo/2015/06/11/07/31/latte-805538_640.jpg'),
(12, 'https://cdn.pixabay.com/photo/2021/01/31/06/36/coffee-5965985_640.jpg'),
(12, 'https://images.unsplash.com/photo-1550948309-0d8983dbdcc3?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTJ8fGxhdHRlfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(12, 'https://images.unsplash.com/photo-1585671582673-66144d229f22?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTF8fGxhdHRlfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(13, 'https://images.unsplash.com/photo-1582785513054-8d1bf9d69c1a?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8TWF0Y2hhJTIwTGF0dGV8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(13, 'https://images.unsplash.com/photo-1515823064-d6e0c04616a7?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8TWF0Y2hhJTIwTGF0dGV8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(13, 'https://images.unsplash.com/photo-1749280447307-31a68eb38673?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8OHx8TWF0Y2hhJTIwTGF0dGV8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(13, 'https://images.unsplash.com/photo-1727850005779-1e24cac382d4?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTJ8fE1hdGNoYSUyMExhdHRlfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(14, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQYpVbR5Su5bzZeDE2o2KMUzESvQqoTIt-khg&s'),
(14, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRQH9bYmAVc4Z42F9Ks0ueLt8LLz_Z6kuy5qg&s'),
(14, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSq_E5RjjWXDs9QBDrLO02pV1A_aE1FQLeSYw&s'),
(14, 'https://images.unsplash.com/photo-1595326202845-d11d23db401b?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NHx8TWlsb3xlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(15, 'https://images.unsplash.com/photo-1597983187002-35a3e1c896d6?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTZ8fENpdmV0JTIwQ29mZmVlfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(15, 'https://images.unsplash.com/photo-1545341122-64b73393fa64?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8N3x8Q2l2ZXQlMjBDb2ZmZWV8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(15, 'https://cdn.pixabay.com/photo/2017/09/04/18/39/coffee-2714970_640.jpg'),
(15, 'https://cdn.pixabay.com/photo/2017/03/13/12/34/coffee-2139592_640.jpg'),
(16, 'https://images.unsplash.com/photo-1629991964567-3107a9202778?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8SWNlZCUyMENvZmZlZSUyMHdpdGglMjBKZWxseXxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(16, 'https://images.unsplash.com/photo-1637378591514-67a2419660fe?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTV8fEljZWQlMjBDb2ZmZWUlMjB3aXRoJTIwSmVsbHl8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(16, 'https://images.unsplash.com/photo-1545136894-70b1ed7cedf6?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTZ8fEljZWQlMjBDb2ZmZWUlMjB3aXRoJTIwSmVsbHl8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(16, 'https://images.unsplash.com/photo-1637378367058-968619b1635c?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTl8fEljZWQlMjBDb2ZmZWUlMjB3aXRoJTIwSmVsbHl8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(17, 'https://cdn.pixabay.com/photo/2022/04/11/08/52/iced-tea-7125271_640.jpg'),
(17, 'https://cdn.pixabay.com/photo/2022/04/11/08/53/iced-tea-7125274_1280.jpg'),
(17, 'https://cdn.pixabay.com/photo/2015/05/03/07/37/tea-750850_640.jpg'),
(17, 'https://cdn.pixabay.com/photo/2016/10/09/17/06/ice-lemon-tea-1726270_640.jpg'),
(18, 'https://cdn.pixabay.com/photo/2016/10/14/18/21/tee-1740871_640.jpg'),
(18, 'https://cdn.pixabay.com/photo/2016/09/21/22/00/tea-1685847_1280.jpg'),
(18, 'https://cdn.pixabay.com/photo/2016/12/19/12/08/ginger-1918107_640.jpg'),
(18, 'https://cdn.pixabay.com/photo/2017/03/01/05/12/tea-cup-2107599_640.jpg'),
(19, 'https://images.unsplash.com/photo-1586195831800-24f14c992cea?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Nnx8Q2hvY29sYXRlJTIwTWlsa3xlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(19, 'https://images.unsplash.com/photo-1542990253-0d0f5be5f0ed?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8Q2hvY29sYXRlJTIwTWlsa3xlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(19, 'https://images.unsplash.com/photo-1517578239113-b03992dcdd25?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTB8fENob2NvbGF0ZSUyME1pbGt8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(19, 'https://images.unsplash.com/photo-1591864384134-8a21ffb51cb5?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTJ8fENob2NvbGF0ZSUyME1pbGt8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(20, 'https://images.unsplash.com/photo-1740103639723-87c39575b2c1?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8VmFuaWxsYSUyME1pbGtzaGFrZXxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(20, 'https://images.unsplash.com/photo-1583152329931-9a3c0c8a91a1?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NHx8VmFuaWxsYSUyME1pbGtzaGFrZXxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(20, 'https://images.unsplash.com/photo-1447195047884-0f014b0d9288?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTB8fFZhbmlsbGElMjBNaWxrc2hha2V8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(20, 'https://images.unsplash.com/photo-1755835070338-6049da75951e?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Nnx8VmFuaWxsYSUyME1pbGtzaGFrZXxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(21, 'https://images.unsplash.com/photo-1685967836529-b0e8d6938227?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8QmFuYW5hJTIwU21vb3RoaWV8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(21, 'https://images.unsplash.com/photo-1707219811295-0f283760668b?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8QmFuYW5hJTIwU21vb3RoaWV8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(21, 'https://images.unsplash.com/photo-1542444592-e6880e286bb9?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NHx8QmFuYW5hJTIwU21vb3RoaWV8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(21, 'https://images.unsplash.com/photo-1542444592-0d5997f202eb?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Nnx8QmFuYW5hJTIwU21vb3RoaWV8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(22, 'https://cdn.pixabay.com/photo/2018/05/07/11/22/mango-3380631_640.jpg'),
(22, 'https://cdn.pixabay.com/photo/2018/05/07/11/22/mango-3380630_640.jpg'),
(22, 'https://plus.unsplash.com/premium_photo-1695055513621-0fecc05fbe38?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTd8fGZyZXNoJTIwbWFuZ28lMjBqdWljZXxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(23, 'https://cdn.pixabay.com/photo/2021/11/05/14/54/tea-6771589_640.jpg'),
(23, 'https://cdn.pixabay.com/photo/2016/11/16/21/35/infused-water-1830120_640.jpg'),
(23, 'https://cdn.pixabay.com/photo/2019/12/19/00/21/ginger-4705132_640.jpg'),
(23, 'https://cdn.pixabay.com/photo/2021/02/01/11/45/food-5970320_640.jpg'),
(24, 'https://cdn.pixabay.com/photo/2023/11/22/20/30/ai-generated-8406387_640.png'),
(24, 'https://cdn.pixabay.com/photo/2017/08/05/22/29/wine-2586196_640.jpg'),
(24, 'https://cdn.pixabay.com/photo/2024/01/23/09/16/ai-generated-8527256_640.jpg'),
(24, 'https://cdn.pixabay.com/photo/2018/05/17/15/33/cocktail-3408834_640.jpg'),
(25, 'https://cdn.pixabay.com/photo/2020/03/07/05/18/beverage-4908765_1280.jpg'),
(25, 'https://images.unsplash.com/photo-1641665271888-575e46923776?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8b3JlbyUyMG1pbGtzaGFrZXxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(25, 'https://images.unsplash.com/photo-1619158401201-8fa932695178?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8b3JlbyUyMG1pbGtzaGFrZXxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(25, 'https://images.unsplash.com/photo-1653852883277-c4b4b9e020e5?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8N3x8b3JlbyUyMG1pbGtzaGFrZXxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(26, 'https://images.unsplash.com/photo-1515316416554-1af1fbca57a2?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8SWNlZCUyMENob2NvbGF0ZXxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(26, 'https://cdn.pixabay.com/photo/2024/08/06/06/52/ai-generated-8948543_640.jpg'),
(27, 'https://images.unsplash.com/photo-1539122144704-71b13a5c0fc7?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Nnx8UHVsbGVkJTIwVGVhfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(28, 'https://images.unsplash.com/photo-1597579018905-8c807adfbed4?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8R3JpbGxlZCUyMENoaWNrZW4lMjBTYW5kd2ljaHxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(28, 'https://images.unsplash.com/photo-1642442928984-b18c2c86b9c2?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8R3JpbGxlZCUyMENoaWNrZW4lMjBTYW5kd2ljaHxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(28, 'https://images.unsplash.com/photo-1757961048411-73703e333d25?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTB8fEdyaWxsZWQlMjBDaGlja2VuJTIwU2FuZHdpY2h8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(28, 'https://images.unsplash.com/photo-1728774283140-7b28a5045502?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8OHx8R3JpbGxlZCUyMENoaWNrZW4lMjBTYW5kd2ljaHxlbnwwfHwwfHx8MA%3D%3D&auto=format&fit=crop&q=60&w=500'),
(29, 'https://cdn.pixabay.com/photo/2022/06/14/18/58/beef-burger-7262651_1280.jpg'),
(29, 'https://cdn.pixabay.com/photo/2020/01/28/12/52/burger-4799767_640.jpg'),
(29, 'https://cdn.pixabay.com/photo/2023/10/25/04/27/burger-8339491_640.jpg'),
(29, 'https://cdn.pixabay.com/photo/2022/07/15/18/17/spicy-burger-7323694_640.jpg'),
(30, 'https://cdn.pixabay.com/photo/2019/04/13/02/39/potato-chips-4123734_640.jpg'),
(30, 'https://cdn.pixabay.com/photo/2017/03/23/04/04/macro-2167281_640.jpg'),
(30, 'https://cdn.pixabay.com/photo/2024/01/23/08/11/ai-generated-8527044_640.jpg'),
(30, 'https://cdn.pixabay.com/photo/2016/11/20/09/06/bowl-1842294_640.jpg'),
(31, 'https://cdn.pixabay.com/photo/2016/05/26/16/27/cinnamon-rolls-1417494_640.jpg'),
(31, 'https://cdn.pixabay.com/photo/2023/11/13/13/19/food-8385524_640.jpg'),
(31, 'https://cdn.pixabay.com/photo/2022/09/11/21/56/pastries-7448102_640.jpg'),
(31, 'https://cdn.pixabay.com/photo/2023/09/09/09/18/cinnamon-rolls-8242739_640.jpg'),
(32, 'https://cdn.pixabay.com/photo/2023/03/22/11/08/ai-generated-7869200_640.jpg'),
(32, 'https://cdn.pixabay.com/photo/2024/08/29/11/37/ai-generated-9006672_640.jpg'),
(32, 'https://cdn.pixabay.com/photo/2015/04/28/14/32/cheesecake-743754_640.jpg'),
(32, 'https://cdn.pixabay.com/photo/2017/09/10/11/56/pastries-2735349_640.jpg'),
(33, 'https://cdn.pixabay.com/photo/2024/11/15/08/40/ai-generated-9198827_640.jpg'),
(33, 'https://cdn.pixabay.com/photo/2022/06/23/14/40/bread-7279975_640.jpg'),
(33, 'https://cdn.pixabay.com/photo/2018/09/03/22/43/croissants-3652582_640.jpg'),
(33, 'https://cdn.pixabay.com/photo/2020/01/05/07/32/breakfast-4742465_640.jpg'),
(34, 'https://cdn.pixabay.com/photo/2017/10/12/17/31/apple-tart-2845186_640.jpg'),
(34, 'https://cdn.pixabay.com/photo/2020/06/19/05/17/egg-tart-5315723_640.jpg'),
(34, 'https://cdn.pixabay.com/photo/2017/08/22/17/55/egg-tart-2669935_1280.jpg'),
(34, 'https://cdn.pixabay.com/photo/2020/03/17/00/39/food-4938704_640.jpg'),
(35, 'https://cdn.pixabay.com/photo/2023/08/17/18/16/ai-generated-8197129_640.jpg'),
(35, 'https://cdn.pixabay.com/photo/2024/09/06/21/44/ai-generated-9028590_640.jpg'),
(35, 'https://cdn.pixabay.com/photo/2023/08/17/19/43/ai-generated-8197292_640.jpg'),
(35, 'https://cdn.pixabay.com/photo/2021/03/31/17/44/donut-6140150_640.jpg'),
(36, 'https://cdn.pixabay.com/photo/2016/07/31/17/51/chicken-1559548_1280.jpg'),
(36, 'https://cdn.pixabay.com/photo/2024/02/02/02/32/ai-generated-8547127_640.jpg'),
(36, 'https://cdn.pixabay.com/photo/2019/02/18/03/10/wings-4003652_640.jpg'),
(36, 'https://cdn.pixabay.com/photo/2019/10/21/17/12/wings-4566612_640.jpg'),
(37, 'https://cdn.pixabay.com/photo/2024/04/29/08/19/ai-generated-8727162_640.jpg'),
(37, 'https://cdn.pixabay.com/photo/2016/08/09/11/24/sandwich-1580348_1280.jpg'),
(37, 'https://cdn.pixabay.com/photo/2023/05/29/17/01/hamburger-8026582_640.jpg'),
(37, 'https://cdn.pixabay.com/photo/2020/06/06/05/31/juicy-sandwich-5265292_640.jpg'),
(38, 'https://cdn.pixabay.com/photo/2024/05/03/15/55/zoni-8737326_640.jpg'),
(38, 'https://cdn.pixabay.com/photo/2015/05/02/01/04/miso-soup-749368_640.jpg'),
(38, 'https://cdn.pixabay.com/photo/2017/10/22/16/04/fish-soup-2878191_640.jpg'),
(38, 'https://cdn.pixabay.com/photo/2023/05/27/13/49/soup-8021570_640.jpg'),
(39, 'https://cdn.pixabay.com/photo/2018/04/24/20/12/dessert-3347905_640.jpg'),
(39, 'https://cdn.pixabay.com/photo/2019/03/18/09/39/brownies-4062740_640.jpg'),
(39, 'https://cdn.pixabay.com/photo/2018/09/19/17/54/chocolate-cake-3689092_640.jpg'),
(39, 'https://cdn.pixabay.com/photo/2024/02/16/16/48/ai-generated-8577991_640.png'),
(40, 'https://cdn.pixabay.com/photo/2017/03/14/19/07/watermelon-2144111_640.jpg'),
(40, 'https://cdn.pixabay.com/photo/2024/09/10/12/38/ai-generated-9037066_640.jpg'),
(40, 'https://cdn.pixabay.com/photo/2024/02/23/12/02/ai-generated-8591923_640.jpg'),
(40, 'https://cdn.pixabay.com/photo/2021/06/13/12/34/fruit-salad-6333173_640.jpg'),
(41, 'https://images.unsplash.com/photo-1648867134727-0b868ba73eb4?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8RXh0cmElMjBFc3ByZXNzbyUyMFNob3R8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(41, 'https://images.unsplash.com/photo-1689164255885-d172a4f5fc4c?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8OHx8RXh0cmElMjBFc3ByZXNzbyUyMFNob3R8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(42, 'https://cdn.pixabay.com/photo/2020/03/07/17/30/cake-4910417_640.jpg'),
(42, 'https://cdn.pixabay.com/photo/2017/08/06/04/16/cupcake-2588646_640.jpg'),
(42, 'https://cdn.pixabay.com/photo/2015/04/02/10/44/cream-703653_640.jpg'),
(42, 'https://cdn.pixabay.com/photo/2017/01/09/11/32/chocolate-hazelnut-1966278_640.jpg'),
(43, 'https://cdn.pixabay.com/photo/2018/04/12/09/17/mat-3312885_640.jpg'),
(44, 'https://cdn.pixabay.com/photo/2016/11/19/15/29/food-1839858_640.jpg'),
(45, 'https://cdn.pixabay.com/photo/2019/02/13/16/28/waffle-3994630_640.jpg'),
(46, 'https://images.unsplash.com/photo-1753095290171-c68a203705f7?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTV8fEV4dHJhJTIwaWNlfGVufDB8fDB8fHww&auto=format&fit=crop&q=60&w=500'),
(47, 'https://cdn.pixabay.com/photo/2019/08/12/20/01/boba-4402053_640.jpg'),
(47, 'https://images.unsplash.com/photo-1734770580735-796a00e42cb2?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8N3x8YnJvd24lMjBzdWdhciUyMGJvYmF8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(47, 'https://images.unsplash.com/photo-1741244133076-afcdda4befae?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Nnx8YnJvd24lMjBzdWdhciUyMGJvYmF8ZW58MHx8MHx8fDA%3D&auto=format&fit=crop&q=60&w=500'),
(48, 'https://cdn.pixabay.com/photo/2016/02/28/21/38/cheese-bread-1227623_640.jpg'),
(49, 'https://cdn.pixabay.com/photo/2025/10/26/16/46/pancakes-9918402_640.png'),
(49, 'https://cdn.pixabay.com/photo/2017/03/13/13/39/pancakes-2139844_640.jpg'),
(49, 'https://cdn.pixabay.com/photo/2017/09/19/08/52/pancake-2764589_640.jpg'),
(49, 'https://cdn.pixabay.com/photo/2016/11/29/04/49/blueberries-1867398_640.jpg'),
(50, 'https://cdn.pixabay.com/photo/2020/08/06/19/39/cookies-5469009_640.jpg'),
(50, 'https://cdn.pixabay.com/photo/2023/11/17/19/07/cookies-8394894_640.jpg'),
(50, 'https://cdn.pixabay.com/photo/2017/06/26/18/14/cookie-dough-2444628_640.jpg'),
(50, 'https://cdn.pixabay.com/photo/2020/03/04/20/35/cookies-4902717_640.jpg');

-- Minuman (Ice & Hot)
INSERT INTO product_variants (product_id, variant_id) VALUES
(1, 1), (1, 2),   -- Espresso Shot
(2, 1), (2, 2),   -- Caramel Macchiato
(3, 1), (3, 2),   -- Mocha Coffee
(4, 1), (4, 2),   -- Iced Coffee Blend
(5, 1), (5, 2),   -- Vanilla Latte
(6, 1), (6, 2),   -- Iced Americano
(7, 1), (7, 2),   -- Affogato
(8, 1), (8, 2),   -- Avocado Coffee
(9, 1), (9, 2),   -- Iced Coffee with Milk
(10,1), (10,2),   -- Traditional Black Coffee
(11,1), (11,2),   -- Cappuccino
(12,1), (12,2),   -- Latte
(13,1), (13,2),   -- Matcha Latte
(14,1), (14,2),   -- Milo Dinosaur
(15,1), (15,2),   -- Civet Coffee
(16,1), (16,2),   -- Iced Coffee with Jelly
(17,1), (17,2),   -- Iced Sweet Tea
(18,1), (18,2),   -- Hot Tea
(19,1), (19,2),   -- Chocolate Milk
(20,1), (20,2),   -- Vanilla Milkshake
(21,1), (21,2),   -- Banana Smoothie
(22,1), (22,2),   -- Fresh Mango Juice
(23,1), (23,2),   -- Honey Lemon Drink
(24,1), (24,2),   -- Iced Orange Juice
(25,1), (25,2),   -- Oreo Milkshake
(26,1), (26,2),   -- Iced Chocolate
(27,1), (27,2),   -- Pulled Tea
(38,1), (38,2),   -- Extra Espresso Shot
(39,1), (39,2),   -- Whipped Cream Topping
(40,1), (40,2),   -- Caramel Syrup
(41,1), (41,2),   -- Vanilla Syrup
(42,1), (42,2),   -- Chocolate Syrup
(43,1), (43,2),   -- Extra Ice
(46,1), (46,2),   -- Hazelnut Syrup
(47,1), (47,2);   -- Chocolate Chips

-- Makanan / Add-on (Food)
INSERT INTO product_variants (product_id, variant_id) VALUES
(28,3),  -- Grilled Chicken Sandwich
(29,3),  -- Beef Burger
(30,3),  -- French Fries
(31,3),  -- Cinnamon Roll
(32,3),  -- Cheesecake
(33,3),  -- Butter Croissant
(34,3),  -- Egg Tart
(35,3),  -- Glazed Donut
(36,3),  -- Spicy Chicken Wings
(37,3),  -- Tuna Sandwich
(38,3),  -- Soup of the Day
(39,3),  -- Chocolate Brownies
(40,3),  -- Mixed Fruit Salad
(44,3),  -- Brown Sugar Boba
(45,3);  -- Cheese Topping

INSERT INTO product_sizes (product_id, size_id) VALUES
(1, 1), (1, 2), (1, 3),
(2, 1), (2, 2), (2, 3),
(3, 1), (3, 2), (3, 3),
(4, 1), (4, 2), (4, 3),
(5, 1), (5, 2), (5, 3),
(6, 1), (6, 2), (6, 3),
(7, 1), (7, 2), (7, 3),
(8, 1), (8, 2), (8, 3),
(9, 1), (9, 2), (9, 3), 
(10, 4), (10, 5),        
(11, 1), (11, 2), (11, 3),
(12, 1), (12, 2), (12, 3),
(13, 1), (13, 2), (13, 3),
(14, 1), (14, 2), (14, 3),
(15, 4), (15, 5),         
(16, 1), (16, 2), (16, 3),
(17, 1), (17, 2), (17, 3),
(18, 1), (18, 2), (18, 3),
(19, 1), (19, 2), (19, 3),
(20, 1), (20, 2), (20, 3),
(21, 1), (21, 2), (21, 3),
(22, 1), (22, 2), (22, 3),
(23, 1), (23, 2), (23, 3),
(24, 1), (24, 2), (24, 3),
(25, 1), (25, 2), (25, 3),
(26, 1), (26, 2), (26, 3),
(27, 1), (27, 2), (27, 3);


INSERT INTO products_categories (product_id, category_id) VALUES
(1, 1), (1, 2), (1, 5),
(2, 1), (2, 2),
(3, 1), (3, 2),
(4, 1), (4, 2),
(5, 1), (5, 2),
(6, 1), (6, 2),
(7, 1), (7, 2), (7, 5),
(8, 1), (8, 2), (8, 5),
(9, 2),
(10, 2),
(11, 2),
(12, 2), (12, 5),
(13, 2), (13, 3),
(14, 3),
(15, 2), (15, 5),
(16, 2),
(17, 3), (17, 5),   
(18, 3),           
(19, 3), (19, 5),   
(20, 3),
(21, 3),            
(22, 3), (22, 5),   
(23, 3),
(24, 3),           
(25, 3), (25, 5),   
(26, 3),
(27, 3), (27, 5),
(28, 4),                   
(29, 4),                   
(30, 4), (30, 5),          
(31, 4), (31, 5),         
(32, 4), (32, 5),         
(33, 4), (33, 5),         
(34, 4), (34, 5),         
(35, 4),                   
(36, 4), (36, 5),      
(37, 4),                   
(38, 4),                   
(39, 4), (39, 5),      
(40, 4),                   
(41, 5),
(42, 5),
(43, 5),
(44, 5),
(45, 5),
(46, 5),
(47, 5),
(48, 5),
(49, 5),
(50, 5);


INSERT INTO payment_methods (name, image) VALUES
('BRI', 'https://cdn3.iconfinder.com/data/icons/banks-in-indonesia-logo-badge/100/BRI-512.png'),
('BCA', 'https://cdn3.iconfinder.com/data/icons/banks-in-indonesia-logo-badge/100/BCA-512.png'),
('MANDIRI', 'https://cdn3.iconfinder.com/data/icons/banks-in-indonesia-logo-badge/100/Mandiri-512.png'),
('BTN', 'https://cdn3.iconfinder.com/data/icons/banks-in-indonesia-logo-badge/100/Bank_BTN-128.png'),
('OVO', 'https://bloguna.com/wp-content/uploads/2025/06/Logo-OVO-Format-PNG-CDR-EPS-SVG-Kualitas-HD-768x615.png'),
('DANA', 'https://bloguna.com/wp-content/uploads/2025/05/Logo-DANA-Format-PNG-CDR-AI-EPS-SVG-768x420.png'),
('PAYPAL', 'https://images.seeklogo.com/logo-png/39/1/paypal-logo-png_seeklogo-390894.png');


INSERT INTO status (name) VALUES 
('Done'),
('Pending'),
('OnProgress'),
('Waiting');

INSERT INTO shippings (name) VALUES 
('DineIn'),
('DoorDelivery'),
('PickUp');

UPDATE shippings
SET additional_price = 10000
WHERE name = 'DoorDelivery';

INSERT INTO promos (title, description, discount, start, "end")
VALUES
('Diskon Awal Bulan', 'Diskon 10% untuk semua produk', 10, '2025-11-01 00:00:00', '2025-11-30 23:59:59'),
('Promo Flash Sale', 'Diskon 20% untuk produk tertentu', 20, '2025-11-15 00:00:00', '2025-11-16 23:59:59');

INSERT INTO recommended_products (product_id, recommended_id) VALUES
(1, 5), (1, 9), (1, 12), (1, 3),
(2, 1), (2, 4), (2, 7), (2, 10),
(3, 14), (3, 12), (3, 5), (3, 8),
(4, 2), (4, 6), (4, 1), (4, 11),
(5, 9), (5, 1), (5, 14), (5, 3),
(6, 1), (6, 2), (6, 10), (6, 4),
(7, 3), (7, 1), (7, 5), (7, 14),
(8, 2), (8, 3), (8, 9), (8, 12),
(9, 12), (9, 1), (9, 3), (9, 10),
(10, 1), (10, 9), (10, 12), (10, 5),
(11, 3), (11, 14), (11, 7), (11, 1),
(12, 9), (12, 3), (12, 8), (12, 14),
(13, 27), (13, 23), (13, 20), (13, 29),
(14, 11), (14, 7), (14, 9), (14, 3),
(15, 10), (15, 11), (15, 13), (15, 17),
(16, 4), (16, 12), (16, 9), (16, 20),
(17, 23), (17, 21), (17, 19), (17, 28),
(18, 21), (18, 22), (18, 23), (18, 27),
(19, 28), (19, 22), (19, 25), (19, 30),
(20, 17), (20, 23), (20, 30), (20, 27),
(21, 27), (21, 23), (21, 25), (21, 30),
(22, 23), (22, 21), (22, 27), (22, 25),
(23, 30), (23, 22), (23, 24), (23, 21),
(24, 18), (24, 21), (24, 25), (24, 30),
(25, 22), (25, 27), (25, 23), (25, 21),
(26, 19), (26, 22), (26, 24), (26, 20),
(27, 25), (27, 23), (27, 21), (27, 30),
(28, 19), (28, 22), (28, 23), (28, 27),
(29, 33), (29, 37), (29, 32), (29, 40),
(30, 27), (30, 24), (30, 25), (30, 23),
(31, 40), (31, 37), (31, 36), (31, 44),
(32, 37), (32, 35), (32, 31), (32, 40),
(33, 29), (33, 37), (33, 35), (33, 40),
(34, 35), (34, 36), (34, 37), (34, 33),
(35, 34), (35, 37), (35, 33), (35, 31),
(36, 44), (36, 37), (36, 35), (36, 33),
(37, 33), (37, 35), (37, 40), (37, 31),
(38, 41), (38, 39), (38, 45), (38, 42),
(39, 42), (39, 38), (39, 41), (39, 45),
(40, 31), (40, 33), (40, 37), (40, 35),
(41, 39), (41, 38), (41, 42), (41, 45),
(42, 45), (42, 39), (42, 41), (42, 38),
(43, 47), (43, 48), (43, 49), (43, 50),
(44, 31), (44, 36), (44, 37), (44, 33),
(45, 38), (45, 42), (45, 39), (45, 41),
(46, 49), (46, 47), (46, 48), (46, 50),
(47, 43), (47, 48), (47, 49), (47, 50),
(48, 50), (48, 47), (48, 49), (48, 43),
(49, 46), (49, 47), (49, 48), (49, 50),
(50, 48), (50, 47), (50, 49), (50, 46);


INSERT INTO users (fullname, email, password, role) VALUES
('admin', 'admin@gmail.com', 'Admin123', 'admin');
