import React, { useEffect, useState } from 'react';

interface Item {
  id: number;
  name: string;
  category: string;
  image_name: string;
};

interface Category {
  id: number;
  name: string;
}

interface ItemResponse {
  items: Item[];
}

interface CategoryResponse {
  categories: Category[];
}

interface CategoryButtonProps {
  category: string;
  onSelectCategory: (category: string) => void;
  isActive: boolean;
}

const CategoryButton: React.FC<CategoryButtonProps> = ({ category, onSelectCategory, isActive }) => (
  <button
    className={`CategoryButton ${isActive ? 'active' : ''}`}
    onClick={() => onSelectCategory(category)}
  >
    {category}
  </button>
);

const server = process.env.REACT_APP_API_URL || 'http://127.0.0.1:9000';
const placeholderImage = process.env.PUBLIC_URL + '/logo192.png';

interface Prop {
  reload?: boolean;
  onLoadCompleted?: () => void;
}

export const ItemList: React.FC<Prop> = ({ reload = true, onLoadCompleted }) => {
  const [items, setItems] = useState<Item[]>([]);
  const [filteredItems, setFilteredItems] = useState<Item[]>([]);
  const [categories, setCategories] = useState<string[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string>('All');

  // Function to retrieve items
  const fetchItems = () => {
    fetch(`${server}/items`, {
      method: 'GET',
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
    })
    .then(response => response.json())
    .then((data: ItemResponse) => {
      console.log('GET success:', data);
      setItems(data.items);
      const categorySet = new Set(data.items.map(item => item.category));
      const categories: string[] = ['All', ...Array.from(categorySet)];
      setCategories(categories);
      setFilteredItems(data.items);
      console.log('GET success:', data.items);
      onLoadCompleted && onLoadCompleted();
      fetchCategories();
    })
    .catch(error => {
      console.error('GET error:', error)
    })
  }
  // Function to retrieve categories
  const fetchCategories = () => {
    fetch(`${server}/categories`, {
      method: 'GET',
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
    })
    .then(response => response.json())
    .then((data: CategoryResponse) => {
      console.log('GET success:', data);
      const categories = ['All', ...data.categories.map(category => category.name)];
      setCategories(categories);
    })
    .catch(error => {
      console.error('GET categories error:', error);
    });
  }

  // Processing when selecting a category
  const handleSelectCategory = (category: string) => {
    setSelectedCategory(category);
    if (category === 'All') {
      console.log('All')
      setFilteredItems(items);
    } else {
      console.log(category)
      setFilteredItems(items.filter(item => item.category === category));
    }
  };

  const ItemDisplay: React.FC<{ item: Item }> = ({ item }) => {
    const imageUrl = item.image_name ? `${server}/image/${item.image_name}` : placeholderImage;

    return (
      <div className="ItemDisplay">
        <img src={imageUrl} alt={item.name} className="ItemImage" />
        <p>{item.name}</p>
      </div>
    );
  };

  useEffect(() => {
    if (reload) {
      fetchItems();
    }
  }, [reload]);

  return (
    <div className='Container'>
      <div className='CategorySection'>
          <div className='CategoryList'>
            <h2>Category</h2>
            {categories.map((category) => (
              <CategoryButton
                key={category}
                category={category}
                onSelectCategory={handleSelectCategory}
                isActive={selectedCategory === category}
              />
            ))}
          </div>
      </div>
      <div className='ItemList'>
        {filteredItems.map((item) => (
          <ItemDisplay key={item.id} item={item} />
        ))}
      </div>
    </div>
  );
};