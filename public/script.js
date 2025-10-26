document.addEventListener('DOMContentLoaded', () => {
    const messageInput = document.getElementById('messageInput');
    const sendButton = document.getElementById('sendButton');
    const messageArea = document.getElementById('messageArea');
    const typingIndicator = document.getElementById('typingIndicator');

    const addMessage = (sender, content) => {
        const messageElement = document.createElement('div');
        messageElement.classList.add('message', sender);

        const senderElement = document.createElement('div');
        senderElement.classList.add('sender');
        senderElement.textContent = sender === 'user' ? 'Anda' : 'GLM-4.5-Flash';

        const contentElement = document.createElement('div');
        
        if (sender === 'ai') {
            contentElement.innerHTML = marked.parse(content);
        } else {
            contentElement.textContent = content;
        }

        messageElement.appendChild(senderElement);
        messageElement.appendChild(contentElement);
        messageArea.appendChild(messageElement);
        
        messageArea.scrollTop = messageArea.scrollHeight;
    };

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

    const sendMessage = async () => {
        const message = messageInput.value.trim();
        if (message === '') return;

        addMessage('user', message);
        messageInput.value = '';

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

            setTypingIndicator(false);

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

    sendButton.addEventListener('click', sendMessage);

    messageInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            sendMessage();
        }
    });

    addMessage('ai', 'Halo! Saya adalah asisten AI yang didukung oleh GLM-4.5-Flash. Ada yang bisa saya bantu?');
});