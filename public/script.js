document.addEventListener('DOMContentLoaded', () => {
    const messageInput = document.getElementById('messageInput');
    const sendButton = document.getElementById('sendButton');
    const messageArea = document.getElementById('messageArea');
    const typingIndicator = document.getElementById('typingIndicator');

    // Fungsi untuk menambahkan pesan ke area chat
    const addMessage = (sender, text) => {
        const messageElement = document.createElement('div');
        messageElement.classList.add('message', sender);

        const senderElement = document.createElement('div');
        senderElement.classList.add('sender');
        senderElement.textContent = sender === 'user' ? 'Anda' : 'GLM-4.5-Flash';

        const textElement = document.createElement('div');
        textElement.textContent = text;

        messageElement.appendChild(senderElement);
        messageElement.appendChild(textElement);
        messageArea.appendChild(messageElement);
        
        // Scroll ke bawah
        messageArea.scrollTop = messageArea.scrollHeight;
    };

    // Fungsi untuk menampilkan/menyembunyikan indikator mengetik
    const setTypingIndicator = (show) => {
        if (show) {
            const indicator = document.createElement('div');
            indicator.id = 'typingIndicator';
            indicator.classList.add('typing-indicator');
            indicator.innerHTML = '<span></span><span></span><span></span>';
            messageArea.appendChild(indicator);
            messageArea.scrollTop = messageArea.scrollHeight;
        } else {
            const indicator = document.getElementById('typingIndicator');
            if (indicator) {
                indicator.remove();
            }
        }
    };

    // Fungsi untuk mengirim pesan
    const sendMessage = async () => {
        const message = messageInput.value.trim();
        if (message === '') return;

        // Tambahkan pesan user ke UI
        addMessage('user', message);
        messageInput.value = '';

        // Tampilkan indikator mengetik
        setTypingIndicator(true);

        try {
            const response = await fetch('/api/chat', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ message: message }),
            });

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            const data = await response.json();

            // Sembunyikan indikator mengetik
            setTypingIndicator(false);

            // Tambahkan pesan AI ke UI
            if (data.reply) {
                addMessage('ai', data.reply);
            } else {
                addMessage('ai', 'Maaf, saya tidak menerima balasan yang valid.');
            }

        } catch (error) {
            console.error('Error:', error);
            setTypingIndicator(false);
            addMessage('ai', 'Maaf, terjadi kesalahan saat menghubungi server.');
        }
    };

    // Event listener untuk tombol kirim
    sendButton.addEventListener('click', sendMessage);

    // Event listener untuk tombol Enter di input
    messageInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            sendMessage();
        }
    });

    // Pesan pembuka dari AI
    addMessage('ai', 'Halo! Saya adalah asisten AI yang didukung oleh GLM-4.5-Flash. Ada yang bisa saya bantu?');
});