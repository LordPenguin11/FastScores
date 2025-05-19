import { Admin, Resource, ListGuesser } from 'react-admin';
import simpleRestProvider from 'ra-data-simple-rest';

const dataProvider = simpleRestProvider('http://localhost:3000');

function App() {
  return (
    <Admin dataProvider={dataProvider}>
      <Resource name="leagues" list={ListGuesser} />
    </Admin>
  );
}

export default App;
