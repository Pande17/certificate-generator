
import { Button, Input } from "antd";


const App = () => {
    const { TextArea } = Input;
  return (
    <div style={{ padding: 20 }}>
      <TextArea placeholder="Masukkan teks" />
      <Button type="primary" style={{ marginTop: 10 }}>
        Klik Saya
      </Button>
    </div>
  );
};

export default App;
